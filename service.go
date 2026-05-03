package main

import (
	"context"
	"errors"
	"sync"
)

const (
	TriggerManual    = "manual"
	TriggerTray      = "tray"
	TriggerScheduler = "scheduler"
)

var ErrBenchmarkInProgress = errors.New("benchmark in progress")

type SiteCheckService struct {
	store             *SettingsStore
	monitor           *ConnectivityMonitor
	telemetry         *TelemetryClient
	onSettingsSaved     func(Settings)
	onBenchmarkFinish   func(BenchmarkReport)
	onDNSFinish         func(DNSTestReport)
	onShowSettings      func()
	onQuit              func()
	dnsManualCtx        context.Context
	dnsManualCancel     context.CancelFunc
	dnsRunningSource    string
	dnsMu               sync.Mutex
	connManualCtx       context.Context
	connManualCancel    context.CancelFunc
	connRunningSource   string
	connMu              sync.Mutex
}

func triggerPriority(source string) int {
	switch source {
	case TriggerManual:
		return 3
	case TriggerTray:
		return 2
	case TriggerScheduler:
		return 1
	default:
		return 0
	}
}

func NewSiteCheckService(store *SettingsStore, monitor *ConnectivityMonitor, telemetry *TelemetryClient) *SiteCheckService {
	return &SiteCheckService{
		store:     store,
		monitor:   monitor,
		telemetry: telemetry,
	}
}

func (s *SiteCheckService) GetSettings() (Settings, error) {
	return s.store.Load()
}

func (s *SiteCheckService) SaveSettings(settings Settings) (Settings, error) {
	current, err := s.store.Load()
	if err != nil {
		current = DefaultSettings()
	}

	saved, err := s.store.Save(settings)
	if err != nil {
		return Settings{}, err
	}
	if current.IntervalMinutes != saved.IntervalMinutes {
		s.telemetry.Track("interval_changed", telemetryProps{
			"interval_minutes": saved.IntervalMinutes,
		})
	}
	if s.onSettingsSaved != nil {
		s.onSettingsSaved(saved)
	}
	return saved, nil
}

func (s *SiteCheckService) Benchmark(source string) (BenchmarkReport, error) {
	s.connMu.Lock()
	if s.connManualCancel != nil {
		if triggerPriority(source) > triggerPriority(s.connRunningSource) || (source == TriggerManual && s.connRunningSource == TriggerManual) {
			s.connManualCancel() // supersedes running one
		} else {
			s.connMu.Unlock()
			return BenchmarkReport{}, ErrBenchmarkInProgress // skip
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.connManualCtx = ctx
	s.connManualCancel = cancel
	s.connRunningSource = source
	s.connMu.Unlock()

	defer func() {
		s.connMu.Lock()
		if s.connManualCtx == ctx { // only clear if it's still ours
			s.connManualCancel = nil
			s.connManualCtx = nil
			s.connRunningSource = ""
		}
		s.connMu.Unlock()
	}()

	settings, err := s.store.Load()
	if err != nil {
		return BenchmarkReport{}, err
	}

	s.telemetry.Track("benchmark_started", telemetryProps{
		"targets_count": len(settings.Targets),
		"source":        source,
	})
	report := s.monitor.Benchmark(ctx, settings.Targets)
	
	if ctx.Err() != nil {
		return BenchmarkReport{}, ctx.Err()
	}

	s.telemetry.Track("benchmark_finished", telemetryProps{
		"targets_count": len(report.Results),
		"has_results":   report.Summary.HasResults,
		"fastest_ms":    report.Summary.FastestMS,
		"slowest_ms":    report.Summary.SlowestMS,
	})
	if s.onBenchmarkFinish != nil {
		s.onBenchmarkFinish(report)
	}
	return report, nil
}

func (s *SiteCheckService) BenchmarkDNS(source string) (DNSTestReport, error) {
	s.dnsMu.Lock()
	if s.dnsManualCancel != nil {
		if triggerPriority(source) > triggerPriority(s.dnsRunningSource) || (source == TriggerManual && s.dnsRunningSource == TriggerManual) {
			s.dnsManualCancel() // supersedes running one
		} else {
			s.dnsMu.Unlock()
			return DNSTestReport{}, ErrBenchmarkInProgress // skip
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.dnsManualCtx = ctx
	s.dnsManualCancel = cancel
	s.dnsRunningSource = source
	s.dnsMu.Unlock()

	defer func() {
		s.dnsMu.Lock()
		if s.dnsManualCtx == ctx { // only clear if it's still ours
			s.dnsManualCancel = nil
			s.dnsManualCtx = nil
			s.dnsRunningSource = ""
		}
		s.dnsMu.Unlock()
	}()

	report := RunDNSTest(ctx)

	if ctx.Err() != nil {
		return DNSTestReport{}, ctx.Err() // canceled
	}

	if s.onDNSFinish != nil {
		s.onDNSFinish(report)
	}
	return report, nil
}


func (s *SiteCheckService) ShowSettings() {
	s.telemetry.Track("settings_opened", nil)
	if s.onShowSettings != nil {
		s.onShowSettings()
	}
}

func (s *SiteCheckService) Quit() {
	if s.onQuit != nil {
		s.onQuit()
	}
}
