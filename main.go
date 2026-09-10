package main

import (
	"log"
	"net/http"

	"GoWebService/internal/handlers"
	"GoWebService/internal/middleware"
)

func main() {
	logger := handlers.NewLogger()

	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/healthcheck", handlers.HealthCheckHandler)
	mux.HandleFunc("/api", handlers.APIHandler)

	// Prometheus metrics endpoint (scraped by Prometheus in the cluster)
	mux.Handle("/metrics", middleware.MetricsHandler())

	// Paths excluded from request logging: the healthcheck is probed
	// constantly by kubelet and the metrics endpoint is scraped regularly;
	// logging them would flood Loki.
	skipLog := map[string]bool{
		"/healthcheck": true,
		"/metrics":     true,
	}

	// Middleware chain (outermost first):
	//   1. request ID  - correlates logs (-> Loki) and future metrics
	//   2. logging     - one structured JSON line per request to stdout,
	//                    collected by fluent-bit and shipped to Loki
	//   3. metrics     - prometheus counters/histograms exposed at /metrics
	var handler http.Handler = mux
	handler = middleware.Metrics(handler)
	handler = middleware.Logging(logger, skipLog)(handler)
	handler = middleware.RequestID(handler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
