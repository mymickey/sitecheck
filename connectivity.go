package main

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

const probeTimeout = 30 * time.Second

type ConnectivityMonitor struct {
	client *http.Client
}

func NewConnectivityMonitor(client *http.Client) *ConnectivityMonitor {
	if client == nil {
		client = &http.Client{Timeout: probeTimeout}
	}
	return &ConnectivityMonitor{client: client}
}

func (m *ConnectivityMonitor) ProbeTarget(ctx context.Context, target Target) ProbeResult {
	started := time.Now()
	result := ProbeResult{
		ID:        target.ID,
		Name:      target.Name,
		URL:       target.URL,
		IconURL:   target.IconURL,
		Status:    StatusUnavailable,
		CheckedAt: started.Format(time.RFC3339),
	}

	probeCtx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(probeCtx, http.MethodGet, target.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	request.Header.Set("Cache-Control", "no-store")
	request.Header.Set("Pragma", "no-cache")

	response, err := m.client.Do(request)
	elapsed := time.Since(started)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1024))

	result.Status = StatusAvailable
	result.LatencyMS = latencyMilliseconds(elapsed)
	return result
}

func (m *ConnectivityMonitor) Benchmark(ctx context.Context, targets []Target) BenchmarkReport {
	results := make([]ProbeResult, len(targets))
	var wg sync.WaitGroup

	for index, target := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(index) * 50 * time.Millisecond)
			results[index] = m.ProbeTarget(ctx, target)
		}()
	}

	wg.Wait()
	return BenchmarkReport{
		Results: results,
		Summary: SummarizeBenchmark(results),
	}
}

func SummarizeBenchmark(results []ProbeResult) BenchmarkSummary {
	summary := BenchmarkSummary{
		CheckedAt: time.Now().Format(time.RFC3339),
	}

	for _, result := range results {
		if result.Status != StatusAvailable || result.LatencyMS <= 0 {
			continue
		}
		if !summary.HasResults {
			summary.FastestMS = result.LatencyMS
			summary.SlowestMS = result.LatencyMS
			summary.HasResults = true
			continue
		}
		if result.LatencyMS < summary.FastestMS {
			summary.FastestMS = result.LatencyMS
		}
		if result.LatencyMS > summary.SlowestMS {
			summary.SlowestMS = result.LatencyMS
		}
	}

	return summary
}

func latencyMilliseconds(duration time.Duration) int {
	milliseconds := int(duration.Milliseconds())
	if milliseconds <= 0 && duration > 0 {
		return 1
	}
	return milliseconds
}
