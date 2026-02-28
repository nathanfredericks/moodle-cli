package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/nathanfredericks/moodle-cli/internal/auth"
	moodleerrors "github.com/nathanfredericks/moodle-cli/internal/errors"
)

// DraftFile represents a file uploaded to Moodle's draft area.
type DraftFile struct {
	ItemID   int    `json:"itemid"`
	FileName string `json:"filename"`
	FileSize int64  `json:"filesize"`
}

// MoodleClient provides a high-level interface to the Moodle API.
type MoodleClient interface {
	Call(ctx context.Context, operationID string, params any, out any) error
	UploadFile(ctx context.Context, filePath string, itemID int) (DraftFile, error)
	BaseURL() string
}

// Client implements MoodleClient.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	retry      moodleerrors.RetryConfig
}

// ClientOptions configures the API client.
type ClientOptions struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
	Retry      moodleerrors.RetryConfig
}

// NewClient creates a new Moodle API client.
func NewClient(opts ClientOptions) *Client {
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	retry := opts.Retry
	if retry.MaxAttempts == 0 {
		retry = moodleerrors.DefaultRetryConfig()
	}
	return &Client{
		baseURL:    opts.BaseURL,
		token:      opts.Token,
		httpClient: httpClient,
		retry:      retry,
	}
}

// NewClientFromConfig creates a client from config and credential store.
func NewClientFromConfig(url string, creds auth.CredentialStore) (*Client, error) {
	token, err := creds.Get(auth.TokenKey)
	if err != nil {
		return nil, &moodleerrors.AuthError{Msg: "no token found; run 'moodle auth login'", Err: err}
	}
	return NewClient(ClientOptions{
		BaseURL: url,
		Token:   token,
	}), nil
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

// Call invokes a Moodle API operation using the standard Moodle REST endpoint.
// Parameters are sent as form-urlencoded POST body with wstoken and wsfunction.
func (c *Client) Call(ctx context.Context, operationID string, params any, out any) error {
	endpoint := fmt.Sprintf("%s/webservice/rest/server.php", c.baseURL)

	// Build form values
	form := url.Values{}
	form.Set("wstoken", c.token)
	form.Set("wsfunction", operationID)
	form.Set("moodlewsrestformat", "json")

	// Flatten params into form values
	if params != nil {
		if err := flattenParams(form, "", params); err != nil {
			return &moodleerrors.ValidationError{Message: fmt.Sprintf("failed to encode params: %v", err)}
		}
	}

	var lastErr error
	for attempt := range c.retry.MaxAttempts {
		body := strings.NewReader(form.Encode())
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
		if err != nil {
			return &moodleerrors.NetworkError{Msg: "failed to create request", Err: err}
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = &moodleerrors.NetworkError{Msg: "request failed", Err: err}
			if attempt < c.retry.MaxAttempts-1 {
				time.Sleep(moodleerrors.RetryDelay(c.retry, attempt, ""))
				continue
			}
			return lastErr
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return &moodleerrors.NetworkError{Msg: "failed to read response", Err: err}
		}

		// Check for retryable status codes
		if moodleerrors.ShouldRetry(resp.StatusCode) && attempt < c.retry.MaxAttempts-1 {
			retryAfter := resp.Header.Get("Retry-After")
			time.Sleep(moodleerrors.RetryDelay(c.retry, attempt, retryAfter))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return &moodleerrors.NetworkError{Msg: fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(respBody))}
		}

		// Check for Moodle error response
		var errResp map[string]any
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if moodleErr := moodleerrors.ParseMoodleError(errResp); moodleErr != nil {
				return moodleErr
			}
		}

		// Unmarshal into output
		if out != nil {
			if err := json.Unmarshal(respBody, out); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

// UploadFile uploads a file to Moodle's draft area via /webservice/upload.php.
// Pass itemID=0 for the first file; use the returned ItemID for subsequent files
// to place them in the same draft area.
func (c *Client) UploadFile(ctx context.Context, filePath string, itemID int) (DraftFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return DraftFile{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("token", c.token)
	_ = writer.WriteField("filearea", "draft")
	_ = writer.WriteField("itemid", fmt.Sprintf("%d", itemID))

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return DraftFile{}, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return DraftFile{}, fmt.Errorf("failed to write file data: %w", err)
	}
	writer.Close()

	endpoint := fmt.Sprintf("%s/webservice/upload.php", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return DraftFile{}, &moodleerrors.NetworkError{Msg: "failed to create upload request", Err: err}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return DraftFile{}, &moodleerrors.NetworkError{Msg: "upload request failed", Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return DraftFile{}, &moodleerrors.NetworkError{Msg: "failed to read upload response", Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return DraftFile{}, &moodleerrors.NetworkError{Msg: fmt.Sprintf("upload failed with status %d: %s", resp.StatusCode, string(respBody))}
	}

	// Check for Moodle error response (single object with "error" key)
	var errResp map[string]any
	if err := json.Unmarshal(respBody, &errResp); err == nil {
		if moodleErr := moodleerrors.ParseMoodleError(errResp); moodleErr != nil {
			return DraftFile{}, moodleErr
		}
	}

	// Response is a JSON array of uploaded file info
	var files []DraftFile
	if err := json.Unmarshal(respBody, &files); err != nil {
		return DraftFile{}, fmt.Errorf("failed to parse upload response: %w", err)
	}
	if len(files) == 0 {
		return DraftFile{}, fmt.Errorf("upload returned empty response")
	}

	return files[0], nil
}

// flattenParams recursively flattens a nested struct/map into form values
// using Moodle's bracket notation (e.g., "courses[0][id]").
func flattenParams(form url.Values, prefix string, v any) error {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			k := fmt.Sprintf("%v", key.Interface())
			newPrefix := k
			if prefix != "" {
				newPrefix = prefix + "[" + k + "]"
			}
			if err := flattenParams(form, newPrefix, rv.MapIndex(key).Interface()); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := range rv.Len() {
			key := fmt.Sprintf("%s[%d]", prefix, i)
			if err := flattenParams(form, key, rv.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Bool:
		if rv.Bool() {
			form.Set(prefix, "1")
		} else {
			form.Set(prefix, "0")
		}
	default:
		form.Set(prefix, fmt.Sprintf("%v", v))
	}
	return nil
}
