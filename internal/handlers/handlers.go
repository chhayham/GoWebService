package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"GoWebService/internal/logic"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func sendJSON(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

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

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("healthcheck requested", "path", r.URL.Path)
	status := logic.HealthStatus()
	sendJSON(w, http.StatusOK, Response{Message: status})
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("api requested", "method", r.Method)
	
	// Input validation: only allow GET
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method)
		sendJSON(w, http.StatusMethodNotAllowed, Response{Error: "Only GET method is allowed"})
		return
	}

	headers := logic.GetHeaderInfo(r)
	sendJSON(w, http.StatusOK, Response{
		Message: "Origin request headers",
		Data:    headers,
	})
}