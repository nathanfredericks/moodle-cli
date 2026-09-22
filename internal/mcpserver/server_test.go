package mcpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()
	NewHandler(Options{APIKey: "secret", Version: "test"}).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
}

func TestMCPRequiresBearerToken(t *testing.T) {
	handler := NewHandler(Options{APIKey: "secret", Version: "test"})
	for _, authorization := range []string{"", "Bearer wrong", "Basic secret"} {
		req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
		req.Header.Set("Authorization", authorization)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusUnauthorized {
			t.Fatalf("authorization %q: status = %d, want %d", authorization, res.Code, http.StatusUnauthorized)
		}
	}
}

func TestMCPInitializeAndListTools(t *testing.T) {
	handler := NewHandler(Options{APIKey: "secret", Version: "test-version"})

	initialize := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`
	initRes := callMCP(t, handler, initialize, "")
	if initRes.Code != http.StatusOK {
		t.Fatalf("initialize status = %d, body = %s", initRes.Code, initRes.Body.String())
	}
	if !strings.Contains(initRes.Body.String(), `"name":"moodle-cli"`) || !strings.Contains(initRes.Body.String(), `"version":"test-version"`) {
		t.Fatalf("initialize body = %s", initRes.Body.String())
	}

	sessionID := initRes.Header().Get("Mcp-Session-Id")
	if sessionID == "" {
		t.Fatal("initialize response did not include Mcp-Session-Id")
	}

	initializedRes := callMCP(t, handler, `{"jsonrpc":"2.0","method":"notifications/initialized"}`, sessionID)
	if initializedRes.Code != http.StatusAccepted {
		t.Fatalf("initialized status = %d, body = %s", initializedRes.Code, initializedRes.Body.String())
	}

	toolsRes := callMCP(t, handler, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`, sessionID)
	if toolsRes.Code != http.StatusOK {
		t.Fatalf("tools/list status = %d, body = %s", toolsRes.Code, toolsRes.Body.String())
	}
	var response struct {
		Result struct {
			Tools []struct {
				Name string `json:"name"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal(toolsRes.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode tools/list: %v; body = %s", err, toolsRes.Body.String())
	}

	names := make(map[string]bool, len(response.Result.Tools))
	for _, tool := range response.Result.Tools {
		names[tool.Name] = true
	}
	for _, required := range []string{"course_list", "assignment_get", "forum_read", "user_whoami", "database_entries"} {
		if !names[required] {
			t.Errorf("tools/list missing %q", required)
		}
	}
}

func TestFlagArgs(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "verbose", value: true, want: "--verbose"},
		{name: "verbose", value: false, want: "--verbose=false"},
		{name: "course", value: float64(42), want: "--course 42"},
		{name: "format", value: "json", want: "--format json"},
	}
	for _, tt := range tests {
		if got := strings.Join(flagArgs(tt.name, tt.value), " "); got != tt.want {
			t.Errorf("flagArgs(%q, %v) = %q, want %q", tt.name, tt.value, got, tt.want)
		}
	}
}

func callMCP(t *testing.T, handler http.Handler, body, sessionID string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
		req.Header.Set("Mcp-Protocol-Version", "2025-11-25")
	}
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res
}
