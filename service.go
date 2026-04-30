package main

import "context"

type SiteCheckService struct {
	store             *SettingsStore
	monitor           *ConnectivityMonitor
	telemetry         *TelemetryClient
	onSettingsSaved   func(Settings)
	onBenchmarkFinish func(BenchmarkReport)
	onShowSettings    func()
	onQuit            func()
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

func (s *SiteCheckService) Benchmark() (BenchmarkReport, error) {
	settings, err := s.store.Load()
	if err != nil {
		return BenchmarkReport{}, err
	}

	s.telemetry.Track("benchmark_started", telemetryProps{
		"targets_count": len(settings.Targets),
	})
	report := s.monitor.Benchmark(context.Background(), settings.Targets)
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
