package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRequestID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := RequestIDFromContext(r.Context())
		if id == "" {
			t.Fatal("expected a request ID in context, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	// The ID should also be echoed back to the client.
	if rec.Header().Get("X-Request-Id") == "" {
		t.Error("expected X-Request-ID to be written to the response")
	}
}

func TestWithRequestID(t *testing.T) {
	// Without an incoming header a new UUID is generated.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	withID := WithRequestID(req)
	if got := RequestIDFromContext(withID.Context()); got == "" {
		t.Fatal("expected WithRequestID to set a request ID")
	}

	// With an incoming header it is preserved.
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("X-Request-ID", "incoming-456")
	withID2 := WithRequestID(req2)
	if got := RequestIDFromContext(withID2.Context()); got != "incoming-456" {
		t.Errorf("expected incoming request ID to be preserved, got %q", got)
	}

	// The original request's context must be untouched.
	if got := RequestIDFromContext(req.Context()); got != "" {
		t.Errorf("expected original request context to have no request ID, got %q", got)
	}
}

func TestRequestIDFromContextWithoutValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := RequestIDFromContext(req.Context()); got != "" {
		t.Errorf("expected empty request ID from a bare context, got %q", got)
	}
}

func TestRequestIDReusesExistingHeader(t *testing.T) {
	var sawID string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawID = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "fixed-id-123")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if sawID != "fixed-id-123" {
		t.Errorf("expected incoming X-Request-Id to be reused, got %q", sawID)
	}
	if rec.Header().Get("X-Request-Id") != "fixed-id-123" {
		t.Errorf("expected response X-Request-Id to be %q, got %q", "fixed-id-123", rec.Header().Get("X-Request-Id"))
	}
}

func TestRequestIDGeneratesValidUUID(t *testing.T) {
	var sawID string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawID = RequestIDFromContext(r.Context())
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	// A UUID v4 is 36 chars: 8-4-4-4-12 hex groups.
	if len(sawID) != 36 {
		t.Fatalf("expected 36-char UUID, got %d chars: %q", len(sawID), sawID)
	}
	if n := strings.Count(sawID, "-"); n != 4 {
		t.Errorf("expected 4 dashes in UUID, got %d: %q", n, sawID)
	}
}

func TestLoggingSkipsExcludedPaths(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	skipLog := map[string]bool{"/healthcheck": true}
	handler := Logging(logger, skipLog)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthcheck", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestLoggingWritesJSONForNormalRequests(t *testing.T) {
	// Capture the log output on a temp file so we can inspect the JSON line.
	f, err := os.CreateTemp(t.TempDir(), "log-*.json")
	if err != nil {
		t.Fatalf("failed to create temp log file: %v", err)
	}
	defer f.Close()

	logger := slog.New(slog.NewJSONHandler(f, nil))
	handler := Logging(logger, map[string]bool{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world")
	}))

	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader("x"))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("failed to rewind log file: %v", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	logged := string(data)

	for _, want := range []string{`"msg":"request"`, `"method":"POST"`, `"path":"/api"`, `"status":200`, `"bytes":11`} {
		if !strings.Contains(logged, want) {
			t.Errorf("expected log line to contain %s, got: %s", want, logged)
		}
	}
}

func TestLoggingRecords404(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "log-*.json")
	if err != nil {
		t.Fatalf("failed to create temp log file: %v", err)
	}
	defer f.Close()

	logger := slog.New(slog.NewJSONHandler(f, nil))
	handler := Logging(logger, map[string]bool{})(http.NotFoundHandler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("failed to rewind log file: %v", err)
	}
	data, _ := io.ReadAll(f)
	if !strings.Contains(string(data), `"status":404`) {
		t.Errorf("expected log line to contain 404, got: %s", data)
	}
}

func TestMetrics(t *testing.T) {
	handler := Metrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		io.WriteString(w, "short")
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusTeapot {
		t.Errorf("expected status 418, got %d", rec.Code)
	}
}

func TestMetricsHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	MetricsHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	// The metrics registered by the metrics package init() must be
	// present in the exposition.
	if !strings.Contains(body, "go_webservice_http_requests_total") {
		t.Errorf("expected body to contain go_webservice_http_requests_total")
	}
	if !strings.Contains(body, "go_webservice_http_requests_in_flight") {
		t.Errorf("expected body to contain go_webservice_http_requests_in_flight")
	}
}
