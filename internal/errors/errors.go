package errors

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// Exit codes
const (
	ExitOK          = 0
	ExitError       = 1
	ExitAuth        = 2
	ExitPermission  = 3
	ExitNotFound    = 4
	ExitValidation  = 5
	ExitNetwork     = 6
	ExitConfig      = 7
	ExitRateLimit   = 8
)

// MoodleError represents an error returned by the Moodle API.
type MoodleError struct {
	ErrorCode string `json:"errorcode"`
	Message   string `json:"message"`
	Exception string `json:"exception"`
	DebugInfo string `json:"debuginfo,omitempty"`
}

func (e *MoodleError) Error() string {
	if e.DebugInfo != "" {
		return fmt.Sprintf("moodle: %s: %s (%s)", e.ErrorCode, e.Message, e.DebugInfo)
	}
	return fmt.Sprintf("moodle: %s: %s", e.ErrorCode, e.Message)
}

// ExitCode maps Moodle error codes to CLI exit codes.
func (e *MoodleError) ExitCode() int {
	switch e.ErrorCode {
	case "invalidtoken", "accessexception":
		return ExitAuth
	case "nopermissions", "requireloginerror":
		return ExitPermission
	case "invalidrecord", "dmlreadingexception":
		return ExitNotFound
	case "invalidparameter", "invalidformdata":
		return ExitValidation
	default:
		return ExitError
	}
}

// AuthError represents authentication failures.
type AuthError struct {
	Msg string
	Err error
}

func (e *AuthError) Error() string { return fmt.Sprintf("auth: %s", e.Msg) }
func (e *AuthError) Unwrap() error { return e.Err }
func (e *AuthError) ExitCode() int { return ExitAuth }

// NetworkError represents connection/transport errors.
type NetworkError struct {
	Msg string
	Err error
}

func (e *NetworkError) Error() string { return fmt.Sprintf("network: %s", e.Msg) }
func (e *NetworkError) Unwrap() error { return e.Err }
func (e *NetworkError) ExitCode() int { return ExitNetwork }

// ValidationError represents invalid user input.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation: %s", e.Message)
}
func (e *ValidationError) ExitCode() int { return ExitValidation }

// ConfigError represents configuration-related errors.
type ConfigError struct {
	Msg string
	Err error
}

func (e *ConfigError) Error() string { return fmt.Sprintf("config: %s", e.Msg) }
func (e *ConfigError) Unwrap() error { return e.Err }
func (e *ConfigError) ExitCode() int { return ExitConfig }

// Exiter is implemented by errors that have an exit code.
type Exiter interface {
	ExitCode() int
}

// GetExitCode returns the exit code for an error.
func GetExitCode(err error) int {
	var exiter Exiter
	if errors.As(err, &exiter) {
		return exiter.ExitCode()
	}
	return ExitError
}

// IsMoodleError checks if the error is a MoodleError with the given code.
func IsMoodleError(err error, code string) bool {
	var me *MoodleError
	if errors.As(err, &me) {
		return me.ErrorCode == code
	}
	return false
}

// ParseMoodleError attempts to parse an error response from the Moodle API.
func ParseMoodleError(body map[string]any) *MoodleError {
	code, _ := body["errorcode"].(string)
	if code == "" {
		return nil
	}
	msg, _ := body["message"].(string)
	exc, _ := body["exception"].(string)
	debug, _ := body["debuginfo"].(string)
	return &MoodleError{
		ErrorCode: code,
		Message:   msg,
		Exception: exc,
		DebugInfo: debug,
	}
}

// RetryConfig configures the retry behavior.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns sensible retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// ShouldRetry determines if an HTTP status code is retryable.
func ShouldRetry(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// RetryDelay computes the delay for the given attempt with jitter.
// If retryAfter is provided from a Retry-After header, it is respected.
func RetryDelay(cfg RetryConfig, attempt int, retryAfter string) time.Duration {
	if retryAfter != "" {
		if secs, err := strconv.Atoi(retryAfter); err == nil {
			return time.Duration(secs) * time.Second
		}
	}

	delay := float64(cfg.BaseDelay) * math.Pow(2, float64(attempt))
	if delay > float64(cfg.MaxDelay) {
		delay = float64(cfg.MaxDelay)
	}
	// Add jitter: ±25%
	jitter := delay * 0.25 * (2*rand.Float64() - 1)
	return time.Duration(delay + jitter)
}
