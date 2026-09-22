package lambdaconfig

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type fakeSecretsManager struct {
	out *secretsmanager.GetSecretValueOutput
	err error
}

func (f fakeSecretsManager) GetSecretValue(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return f.out, f.err
}

func TestLoad(t *testing.T) {
	secret := `{"MOODLE_URL":"https://moodle.example.com","MOODLE_TOKEN":"token","MCP_API_KEY":"key"}`
	cfg, err := Load(context.Background(), fakeSecretsManager{
		out: &secretsmanager.GetSecretValueOutput{SecretString: &secret},
	}, "moodle-mcp")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.MoodleURL != "https://moodle.example.com" || cfg.MoodleToken != "token" || cfg.MCPAPIKey != "key" {
		t.Fatalf("Load() = %#v", cfg)
	}
}

func TestLoadFailures(t *testing.T) {
	valid := `{"MOODLE_URL":"https://moodle.example.com","MOODLE_TOKEN":"token","MCP_API_KEY":"key"}`
	tests := []struct {
		name     string
		secretID string
		client   fakeSecretsManager
		want     string
	}{
		{name: "missing secret id", client: fakeSecretsManager{}, want: "secret ID is required"},
		{name: "AWS error", secretID: "moodle-mcp", client: fakeSecretsManager{err: errors.New("denied")}, want: "load Moodle MCP secret"},
		{name: "binary secret", secretID: "moodle-mcp", client: fakeSecretsManager{out: &secretsmanager.GetSecretValueOutput{}}, want: "JSON SecretString"},
		{name: "invalid JSON", secretID: "moodle-mcp", client: fakeSecretsManager{out: output(`{`)}, want: "decode Moodle MCP secret"},
		{name: "unknown field", secretID: "moodle-mcp", client: fakeSecretsManager{out: output(strings.TrimSuffix(valid, "}") + `,"EXTRA":"value"}`)}, want: "unknown field"},
		{name: "missing URL", secretID: "moodle-mcp", client: fakeSecretsManager{out: output(`{"MOODLE_TOKEN":"token","MCP_API_KEY":"key"}`)}, want: "MOODLE_URL is required"},
		{name: "invalid URL", secretID: "moodle-mcp", client: fakeSecretsManager{out: output(`{"MOODLE_URL":"moodle.example.com","MOODLE_TOKEN":"token","MCP_API_KEY":"key"}`)}, want: "absolute HTTP(S) URL"},
		{name: "missing token", secretID: "moodle-mcp", client: fakeSecretsManager{out: output(`{"MOODLE_URL":"https://moodle.example.com","MCP_API_KEY":"key"}`)}, want: "MOODLE_TOKEN is required"},
		{name: "missing API key", secretID: "moodle-mcp", client: fakeSecretsManager{out: output(`{"MOODLE_URL":"https://moodle.example.com","MOODLE_TOKEN":"token"}`)}, want: "MCP_API_KEY is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(context.Background(), tt.client, tt.secretID)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Load() error = %v, want substring %q", err, tt.want)
			}
		})
	}
}

func TestApplyEnvironment(t *testing.T) {
	t.Setenv("MOODLE_URL", "")
	t.Setenv("MOODLE_TOKEN", "")
	t.Setenv("MCP_API_KEY", "")

	cfg := RuntimeConfig{
		MoodleURL:   "https://moodle.example.com",
		MoodleToken: "token",
		MCPAPIKey:   "key",
	}
	if err := cfg.ApplyEnvironment(); err != nil {
		t.Fatalf("ApplyEnvironment() error = %v", err)
	}
	if got := []string{
		os.Getenv("MOODLE_URL"),
		os.Getenv("MOODLE_TOKEN"),
		os.Getenv("MCP_API_KEY"),
	}; strings.Join(got, ",") != "https://moodle.example.com,token,key" {
		t.Fatalf("environment = %v", got)
	}
}

func output(value string) *secretsmanager.GetSecretValueOutput {
	return &secretsmanager.GetSecretValueOutput{SecretString: &value}
}
