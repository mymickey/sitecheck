package main

import "context"

type SiteCheckService struct {
	store             *SettingsStore
	monitor           *ConnectivityMonitor
	onSettingsSaved   func(Settings)
	onBenchmarkFinish func(BenchmarkReport)
	onShowSettings    func()
	onQuit            func()
}

func NewSiteCheckService(store *SettingsStore, monitor *ConnectivityMonitor) *SiteCheckService {
	return &SiteCheckService{
		store:   store,
		monitor: monitor,
	}
}

func (s *SiteCheckService) GetSettings() (Settings, error) {
	return s.store.Load()
}

func (s *SiteCheckService) SaveSettings(settings Settings) (Settings, error) {
	saved, err := s.store.Save(settings)
	if err != nil {
		return Settings{}, err
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

	report := s.monitor.Benchmark(context.Background(), settings.Targets)
	if s.onBenchmarkFinish != nil {
		s.onBenchmarkFinish(report)
	}
	return report, nil
}

func (s *SiteCheckService) ShowSettings() {
	if s.onShowSettings != nil {
		s.onShowSettings()
	}
}

func (s *SiteCheckService) Quit() {
	if s.onQuit != nil {
		s.onQuit()
	}
}
