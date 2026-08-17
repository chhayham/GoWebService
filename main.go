package main

import (
	"log"
	"net/http"

	"GoWebService/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/healthcheck", handlers.HealthCheckHandler)
	mux.HandleFunc("/api", handlers.APIHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
