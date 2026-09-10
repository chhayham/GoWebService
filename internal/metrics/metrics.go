// Package metrics defines the prometheus collectors exposed at /metrics.
//
// The cluster's Prometheus scrapes the /metrics endpoint; all collectors
// are registered against the default registry, so promhttp.Handler()
// picks them up automatically.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Labels deliberately kept to low-cardinality values (HTTP method and a
// status class like "2xx"). Using the raw URL path as a label would let
// clients create unbounded label values.
var (
	ReqCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "go_webservice",
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests processed, partitioned by method and status class.",
	}, []string{"method", "status"})

	ReqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "go_webservice",
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request latency in seconds, partitioned by method and status class.",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"method", "status"})

	ReqBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "go_webservice",
		Name:      "http_response_size_bytes_total",
		Help:      "Total size of HTTP responses in bytes, partitioned by method and status class.",
	}, []string{"method", "status"})

	ReqInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "go_webservice",
		Name:      "http_requests_in_flight",
		Help:      "Number of HTTP requests currently being served.",
	})
)

func init() {
	prometheus.MustRegister(ReqCount, ReqDur, ReqBytes, ReqInFlight)
}
