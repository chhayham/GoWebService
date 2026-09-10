package logic

import (
	"net/http"
)

// safeHeaders is the allowlist of headers GetHeaderInfo will expose.
// Using an allowlist (instead of a blocklist) means new or sensitive
// headers (Cookie, Authorization, X-Forwarded-*, etc.) are never
// returned unless explicitly added here.
var safeHeaders = map[string]bool{
	"Accept":          true,
	"Accept-Encoding": true,
	"Accept-Language": true,
	"Cache-Control":   true,
	"Content-Type":    true,
	"Origin":          true,
	"Referer":         true,
	"User-Agent":      true,
	"X-Request-Id":    true,
}

// GetHeaderInfo extracts the allowlisted headers from the request and
// returns them as a new map. Sensitive headers are always excluded.
func GetHeaderInfo(r *http.Request) map[string][]string {
	safe := make(map[string][]string)
	for key, values := range r.Header {
		if !safeHeaders[key] {
			continue
		}
		safe[key] = values
	}
	return safe
}

// HealthStatus returns a simple status message
func HealthStatus() string {
	return "OK"
}
