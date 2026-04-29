package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProbeTargetTreatsAnyHTTPStatusAsReachable(t *testing.T) {
	var method string
	var cacheControl string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		cacheControl = r.Header.Get("Cache-Control")
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(server.Close)

	monitor := NewConnectivityMonitor(&http.Client{Timeout: time.Second})
	result := monitor.ProbeTarget(context.Background(), Target{
		ID:   "github",
		Name: "GitHub",
		URL:  server.URL,
	})

	if result.Status != StatusAvailable {
		t.Fatalf("ProbeTarget() status = %q, want %q", result.Status, StatusAvailable)
	}
	if method != http.MethodGet {
		t.Errorf("ProbeTarget() method = %q, want %q", method, http.MethodGet)
	}
	if cacheControl != "no-store" {
		t.Errorf("ProbeTarget() Cache-Control = %q, want no-store", cacheControl)
	}
	if result.LatencyMS <= 0 {
		t.Errorf("ProbeTarget() LatencyMS = %d, want positive latency", result.LatencyMS)
	}
}

func TestSummarizeBenchmarkIgnoresUnavailableTargets(t *testing.T) {
	report := SummarizeBenchmark([]ProbeResult{
		{ID: "fast", Status: StatusAvailable, LatencyMS: 22},
		{ID: "down", Status: StatusUnavailable, LatencyMS: 0},
		{ID: "slow", Status: StatusAvailable, LatencyMS: 128},
	})

	if !report.HasResults {
		t.Fatalf("SummarizeBenchmark() HasResults = false, want true")
	}
	if report.FastestMS != 22 {
		t.Errorf("SummarizeBenchmark() FastestMS = %d, want 22", report.FastestMS)
	}
	if report.SlowestMS != 128 {
		t.Errorf("SummarizeBenchmark() SlowestMS = %d, want 128", report.SlowestMS)
	}
}

func TestSummarizeBenchmarkReportsNoResultsWhenAllUnavailable(t *testing.T) {
	report := SummarizeBenchmark([]ProbeResult{
		{ID: "down", Status: StatusUnavailable, LatencyMS: 0},
	})

	if report.HasResults {
		t.Fatalf("SummarizeBenchmark() HasResults = true, want false")
	}
	if report.FastestMS != 0 || report.SlowestMS != 0 {
		t.Errorf("SummarizeBenchmark() = %+v, want zero latency summary", report)
	}
}
