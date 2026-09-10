package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"GoWebService/internal/metrics"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// statusRecorder wraps ResponseWriter to capture the status code and bytes
// written, so they can be attached to logs and metrics after the handler
// has finished.
type statusRecorder struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (sr *statusRecorder) WriteHeader(status int) {
	if sr.status == 0 {
		sr.status = status
	}
	sr.ResponseWriter.WriteHeader(status)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK
	}
	n, err := sr.ResponseWriter.Write(b)
	sr.bytesWritten += n
	return n, err
}

func (sr *statusRecorder) Flush() {
	if f, ok := sr.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// requestIDKey is the context key used to carry the request ID.
type requestIDKey struct{}

// WithRequestID returns a copy of the request with a request ID set in its
// context. The ID is taken from the X-Request-ID header if present, otherwise
// a generated UUID.
func WithRequestID(r *http.Request) *http.Request {
	id := r.Header.Get("X-Request-ID")
	if id == "" {
		id = uuid.NewString()
	}
	ctx := context.WithValue(r.Context(), requestIDKey{}, id)
	return r.WithContext(ctx)
}

// RequestIDFromContext extracts the request ID from the request context,
// returning an empty string when none is set.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// RequestID is middleware that ensures every request carries an ID (from
// the X-Request-ID header if the client sent one, otherwise a fresh UUID).
// The ID is echoed back in the response so a client can reference it when
// reporting an issue; it also correlates log lines shipped to Loki.
//
// It must be applied outermost (before Logging) so the ID is available to
// downstream middleware and handlers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = WithRequestID(r)
		w.Header().Set("X-Request-ID", RequestIDFromContext(r.Context()))
		next.ServeHTTP(w, r)
	})
}

// Logging is middleware that emits one structured JSON log line per request.
// The logger writes to stdout, where fluent-bit in the cluster collects it
// and pushes it to Loki.
//
// skip is a map of path prefixes that are excluded from logging (e.g. the
// healthcheck endpoint, which is probed frequently and would otherwise
// flood Loki).
func Logging(logger *slog.Logger, skip map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for prefix := range skip {
				if strings.HasPrefix(r.URL.Path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			logger.Info("request",
				"request_id", RequestIDFromContext(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"bytes", rec.bytesWritten,
				"duration_ms", float64(time.Since(start).Microseconds())/1000.0,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}

// Metrics is middleware that records HTTP request duration and in-flight
// request counters via prometheus client_golang. The cluster's Prometheus
// scrapes these from the /metrics endpoint.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		statusClass := strconv.Itoa(rec.status/100) + "xx"

		metrics.ReqInFlight.Inc()
		defer metrics.ReqInFlight.Dec()

		metrics.ReqCount.WithLabelValues(r.Method, statusClass).Inc()
		metrics.ReqDur.WithLabelValues(r.Method, statusClass).Observe(duration.Seconds())
		metrics.ReqBytes.WithLabelValues(r.Method, statusClass).Add(float64(rec.bytesWritten))
	})
}

// MetricsHandler serves the prometheus metrics endpoint. It should be
// wrapped by the same middleware chain as the rest of the app so its
// traffic is logged too.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
