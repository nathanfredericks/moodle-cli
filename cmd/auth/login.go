package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/nathanfredericks/moodle-cli/internal/auth"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/config"
)

type tokenResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

func newLoginCmd(f *cmdutil.Factory) *cobra.Command {
	var moodleURL string
	var username string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to a Moodle instance",
		Long:  "Authenticate with a Moodle instance using username and password to obtain an API token.",
		Example: `  # Log in interactively (prompts for URL, username, and password)
  moodle auth login

  # Log in with URL and username flags
  moodle auth login --url https://moodle.example.com --username jdoe`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Prompt for URL if not provided
			if moodleURL == "" {
				fmt.Fprint(f.IO.Out, "Moodle URL: ")
				var input string
				fmt.Fscanln(f.IO.In, &input)
				moodleURL = strings.TrimSpace(input)
			}
			if moodleURL == "" {
				return fmt.Errorf("URL is required")
			}
			moodleURL = strings.TrimRight(moodleURL, "/")

			// Prompt for username if not provided
			if username == "" {
				fmt.Fprint(f.IO.Out, "Username: ")
				var input string
				fmt.Fscanln(f.IO.In, &input)
				username = strings.TrimSpace(input)
			}
			if username == "" {
				return fmt.Errorf("username is required")
			}

			// Prompt for password (no echo)
			fmt.Fprint(f.IO.Out, "Password: ")
			passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Fprintln(f.IO.Out) // newline after password
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			password := string(passwordBytes)
			if password == "" {
				return fmt.Errorf("password is required")
			}

			// POST to login/token.php
			formData := url.Values{
				"username": {username},
				"password": {password},
				"service":  {"moodle_mobile_app"},
			}
			resp, err := http.PostForm(moodleURL+"/login/token.php", formData)
			if err != nil {
				return fmt.Errorf("failed to connect to %s: %w", moodleURL, err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			var tokenResp tokenResponse
			if err := json.Unmarshal(body, &tokenResp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			if tokenResp.Error != "" {
				return fmt.Errorf("login failed: %s", tokenResp.Error)
			}
			if tokenResp.Token == "" {
				return fmt.Errorf("login failed: no token in response")
			}

			// Store token
			if err := f.Auth.Set(auth.TokenKey, tokenResp.Token); err != nil {
				return fmt.Errorf("failed to save token: %w", err)
			}

			// Save site
			site := config.Site{
				URL:      moodleURL,
				Username: username,
			}
			if err := f.Config.SaveSite(site); err != nil {
				return fmt.Errorf("failed to save site: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Logged in to %s as %s\n", moodleURL, username)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moodleURL, "url", "u", os.Getenv("MOODLE_URL"), "Moodle instance URL")
	cmd.Flags().StringVar(&username, "username", "", "Username")

	return cmd
}
