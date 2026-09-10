package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"GoWebService/internal/logic"
)

// NewLogger builds the application logger: structured JSON to stdout.
// In the cluster, fluent-bit tail-logs the container's stdout and pushes
// the lines to Loki, so keeping the JSON on stdout is what makes the
// logs queryable there.
func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

var logger = NewLogger()

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func sendJSON(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// HomeHandler is registered at "/" on the default mux, which makes it the
// catch-all for every unregistered path. The path check below is therefore
// required: it turns unknown paths into 404s instead of serving the page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<link rel="stylesheet" href="/static/style.css">
		</head>
		<body>
			<div class="container">
				<h1>Hello World</h1>
			</div>
		</body>
		</html>
	`)
}

// HealthCheckHandler reports the service status. This is the endpoint hit
// by the k8s liveness/readiness/startup probes (see helm-chart), which
// fire every 10s per pod, so it must stay lightweight and log nothing —
// the middleware also skips logging for this path.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := logic.HealthStatus()
	sendJSON(w, http.StatusOK, Response{Message: status})
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("api requested", "method", r.Method)

	// Input validation: only allow GET
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method)
		w.Header().Set("Allow", http.MethodGet)
		sendJSON(w, http.StatusMethodNotAllowed, Response{Error: "Only GET method is allowed"})
		return
	}

	headers := logic.GetHeaderInfo(r)
	sendJSON(w, http.StatusOK, Response{
		Message: "Origin request headers",
		Data:    headers,
	})
}
