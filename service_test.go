package main

import (
	"context"
	"testing"
	"time"
)

func TestSiteCheckService_BenchmarkDNS_Priority(t *testing.T) {
	service := NewSiteCheckService(
		NewSettingsStore("test-settings.json"),
		NewConnectivityMonitor(nil),
		NewTelemetryClient("test"),
	)

	// 1. Simulate an in-flight manual run
	service.dnsMu.Lock()
	_, cancel1 := context.WithCancel(context.Background())
	service.dnsManualCancel = cancel1
	service.dnsRunningSource = TriggerManual
	service.dnsMu.Unlock()

	// Tray trigger should be skipped when Manual is running
	_, err := service.BenchmarkDNS(TriggerTray)
	if err != ErrBenchmarkInProgress {
		t.Errorf("Expected ErrBenchmarkInProgress for tray when manual is running, got %v", err)
	}

	// 2. Simulate an in-flight scheduler run
	service.dnsMu.Lock()
	ctx2, cancel2 := context.WithCancel(context.Background())
	service.dnsManualCancel = cancel2
	service.dnsRunningSource = TriggerScheduler
	service.dnsMu.Unlock()

	// Tray trigger SHOULD preempt scheduler run
	go func() {
		_, _ = service.BenchmarkDNS(TriggerTray)
	}()

	select {
	case <-ctx2.Done():
		// Canceled successfully by Tray
	case <-time.After(1 * time.Second):
		t.Errorf("Expected previous scheduler run to be canceled by new tray run")
	}

	// 3. Simulate an in-flight tray run
	service.dnsMu.Lock()
	ctx3, cancel3 := context.WithCancel(context.Background())
	service.dnsManualCancel = cancel3
	service.dnsRunningSource = TriggerTray
	service.dnsMu.Unlock()

	// Manual trigger SHOULD preempt tray run
	go func() {
		_, _ = service.BenchmarkDNS(TriggerManual)
	}()

	select {
	case <-ctx3.Done():
		// Canceled successfully by Manual
	case <-time.After(1 * time.Second):
		t.Errorf("Expected previous tray run to be canceled by new manual run")
	}
}

func TestSiteCheckService_Benchmark_Priority(t *testing.T) {
	service := NewSiteCheckService(
		NewSettingsStore("test-settings.json"),
		NewConnectivityMonitor(nil),
		NewTelemetryClient("test"),
	)

	// 1. Simulate an in-flight manual run
	service.connMu.Lock()
	_, cancel1 := context.WithCancel(context.Background())
	service.connManualCancel = cancel1
	service.connRunningSource = TriggerManual
	service.connMu.Unlock()

	// Tray trigger should be skipped when Manual is running
	_, err := service.Benchmark(TriggerTray)
	if err != ErrBenchmarkInProgress {
		t.Errorf("Expected ErrBenchmarkInProgress for tray when manual is running, got %v", err)
	}

	// 2. Simulate an in-flight scheduler run
	service.connMu.Lock()
	ctx2, cancel2 := context.WithCancel(context.Background())
	service.connManualCancel = cancel2
	service.connRunningSource = TriggerScheduler
	service.connMu.Unlock()

	// Tray trigger SHOULD preempt scheduler run
	go func() {
		_, _ = service.Benchmark(TriggerTray)
	}()

	select {
	case <-ctx2.Done():
		// Canceled successfully by Tray
	case <-time.After(1 * time.Second):
		t.Errorf("Expected previous scheduler run to be canceled by new tray run")
	}

	// 3. Simulate an in-flight tray run
	service.connMu.Lock()
	ctx3, cancel3 := context.WithCancel(context.Background())
	service.connManualCancel = cancel3
	service.connRunningSource = TriggerTray
	service.connMu.Unlock()

	// Manual trigger SHOULD preempt tray run
	go func() {
		_, _ = service.Benchmark(TriggerManual)
	}()

	select {
	case <-ctx3.Done():
		// Canceled successfully by Manual
	case <-time.After(1 * time.Second):
		t.Errorf("Expected previous tray run to be canceled by new manual run")
	}
}
