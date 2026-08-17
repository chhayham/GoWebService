package logic

import (
	"net/http"
)

// GetHeaderInfo extracts headers from the request to be returned as a map
func GetHeaderInfo(r *http.Request) map[string][]string {
	return r.Header
}

// HealthStatus returns a simple status message
func HealthStatus() string {
	return "OK"
}