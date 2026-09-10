package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// TestCollectorsRegistered verifies that all collectors are registered
// against the default registry (in init()), which is what promhttp.Handler()
// and therefore the /metrics endpoint serve.
func TestCollectorsRegistered(t *testing.T) {
	// Empty *Vec collectors emit no series on Gather(); touch each one with
	// a test-only label pair so its metric family shows up in the output.
	_ = ReqCount.WithLabelValues("TEST", "4xx")
	_ = ReqDur.WithLabelValues("TEST", "4xx")
	_ = ReqBytes.WithLabelValues("TEST", "4xx")

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics from default registry: %v", err)
	}

	want := map[string]bool{
		"go_webservice_http_requests_total":            false,
		"go_webservice_http_request_duration_seconds":  false,
		"go_webservice_http_response_size_bytes_total": false,
		"go_webservice_http_requests_in_flight":        false,
	}
	for _, mf := range families {
		if _, ok := want[mf.GetName()]; ok {
			want[mf.GetName()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected metric family %q to be registered on the default registry", name)
		}
	}
}

func TestCounterAndGaugeOperations(t *testing.T) {
	// Use an unusual-but-valid label pair so this test does not depend on
	// (and does not pollute the assertions of) the middleware tests.
	const method, status = "CONNECT", "5xx"

	before := counterValue(t, method, status)

	ReqCount.WithLabelValues(method, status).Inc()
	ReqBytes.WithLabelValues(method, status).Add(42)
	ReqInFlight.Inc()
	defer ReqInFlight.Dec()

	if after := counterValue(t, method, status); after != before+1 {
		t.Errorf("expected request counter to increment from %v to %v, got %v", before, before+1, after)
	}
}

// counterValue reads the current value of the go_webservice_http_requests_total
// counter for a given method/status label pair.
func counterValue(t *testing.T, method, status string) float64 {
	t.Helper()

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}
	for _, mf := range families {
		if mf.GetName() != "go_webservice_http_requests_total" {
			continue
		}
		for _, m := range mf.GetMetric() {
			labels := map[string]string{}
			for _, lp := range m.GetLabel() {
				labels[lp.GetName()] = lp.GetValue()
			}
			if labels["method"] == method && labels["status"] == status {
				return m.GetCounter().GetValue()
			}
		}
	}
	return 0
}
