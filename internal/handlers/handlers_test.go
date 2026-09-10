package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandlerRoot(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("expected text/html content type, got %q", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Hello World") {
		t.Errorf("expected body to contain %q, got: %s", "Hello World", body)
	}
}

func TestHomeHandlerUnknownPathIs404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/definitely-not-registered", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for unknown path, got %d", rec.Code)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rec := httptest.NewRecorder()

	HealthCheckHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json content type, got %q", ct)
	}

	var payload Response
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if payload.Message != "OK" {
		t.Errorf("expected message %q, got %q", "OK", payload.Message)
	}
	if payload.Error != "" {
		t.Errorf("expected no error in payload, got %q", payload.Error)
	}
}

func TestAPIHandlerGetReturnsSafeHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Cookie", "session=secret")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	APIHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload struct {
		Message string              `json:"message"`
		Data    map[string][]string `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	for _, want := range []string{"Content-Type", "User-Agent"} {
		if len(payload.Data[want]) == 0 {
			t.Errorf("expected allowlisted header %q in response data", want)
		}
	}
	for _, banned := range []string{"Cookie", "Authorization"} {
		if values, ok := payload.Data[banned]; ok {
			t.Errorf("expected sensitive header %q to be excluded, got %v", banned, values)
		}
	}
}

func TestAPIHandlerRejectsNonGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api", nil)
	rec := httptest.NewRecorder()

	APIHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != http.MethodGet {
		t.Errorf("expected Allow header %q, got %q", http.MethodGet, allow)
	}

	var payload Response
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if payload.Error == "" {
		t.Error("expected an error message in the 405 response")
	}
}

func TestNewLoggerWritesJSONToStdout(t *testing.T) {
	// Just ensure construction succeeds and it is usable; it writes to
	// stdout (collected by fluent-bit -> Loki in the cluster).
	logger := NewLogger()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
	logger.Info("test", "key", "value")
}
