package logic

import (
	"net/http/httptest"
	"testing"
)

func TestGetHeaderInfo(t *testing.T) {
	// Create a test request with some headers
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")

	// Call the function
	headers := GetHeaderInfo(req)

	// Verify the headers are correctly returned
	if headers == nil {
		t.Errorf("Expected headers to not be nil")
	}

	// Check if specific headers are present
	if contentType := headers["Content-Type"]; len(contentType) == 0 {
		t.Errorf("Expected Content-Type header to be present")
	}

	if userAgent := headers["User-Agent"]; len(userAgent) == 0 {
		t.Errorf("Expected User-Agent header to be present")
	}
}

func TestHealthStatus(t *testing.T) {
	// Call the function
	status := HealthStatus()

	// Verify the status is as expected
	expected := "OK"
	if status != expected {
		t.Errorf("Expected status %s, got %s", expected, status)
	}
}