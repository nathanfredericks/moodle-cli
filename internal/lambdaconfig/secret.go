// Package lambdaconfig loads the fixed Moodle MCP runtime credentials from
// AWS Secrets Manager.
package lambdaconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerAPI is the subset of the AWS client used by the loader.
type SecretsManagerAPI interface {
	GetSecretValue(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// RuntimeConfig is the JSON contract stored in Secrets Manager.
type RuntimeConfig struct {
	MoodleURL   string `json:"MOODLE_URL"`
	MoodleToken string `json:"MOODLE_TOKEN"`
	MCPAPIKey   string `json:"MCP_API_KEY"`
}

// Load retrieves, decodes, and validates the named Secrets Manager secret.
func Load(ctx context.Context, client SecretsManagerAPI, secretID string) (RuntimeConfig, error) {
	if strings.TrimSpace(secretID) == "" {
		return RuntimeConfig{}, fmt.Errorf("secret ID is required")
	}

	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &secretID})
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("load Moodle MCP secret: %w", err)
	}
	if out.SecretString == nil {
		return RuntimeConfig{}, fmt.Errorf("Moodle MCP secret must contain a JSON SecretString")
	}

	var cfg RuntimeConfig
	decoder := json.NewDecoder(strings.NewReader(*out.SecretString))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return RuntimeConfig{}, fmt.Errorf("decode Moodle MCP secret: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return RuntimeConfig{}, err
	}
	return cfg, nil
}

// Validate checks that every required value is present and the Moodle URL is
// an absolute HTTP(S) endpoint.
func (c RuntimeConfig) Validate() error {
	if strings.TrimSpace(c.MoodleURL) == "" {
		return fmt.Errorf("MOODLE_URL is required")
	}
	parsed, err := url.Parse(c.MoodleURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("MOODLE_URL must be an absolute HTTP(S) URL")
	}
	if strings.TrimSpace(c.MoodleToken) == "" {
		return fmt.Errorf("MOODLE_TOKEN is required")
	}
	if strings.TrimSpace(c.MCPAPIKey) == "" {
		return fmt.Errorf("MCP_API_KEY is required")
	}
	return nil
}

// ApplyEnvironment exposes the credentials through the same environment
// contract used by the existing CLI factory.
func (c RuntimeConfig) ApplyEnvironment() error {
	for key, value := range map[string]string{
		"MOODLE_URL":   c.MoodleURL,
		"MOODLE_TOKEN": c.MoodleToken,
		"MCP_API_KEY":  c.MCPAPIKey,
	} {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
	}
	return nil
}
