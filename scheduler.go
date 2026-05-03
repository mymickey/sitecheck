package main

import (
	"sync"
	"time"
)

type BenchmarkScheduler struct {
	service *SiteCheckService
	mu      sync.Mutex
	stop    chan struct{}
}

func NewBenchmarkScheduler(service *SiteCheckService) *BenchmarkScheduler {
	return &BenchmarkScheduler{service: service}
}

func (s *BenchmarkScheduler) Start(intervalMinutes int, dnsIntervalHours int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopLocked()
	stop := make(chan struct{})
	s.stop = stop
	interval := time.Duration(intervalMinutes) * time.Minute
	dnsInterval := time.Duration(dnsIntervalHours) * time.Hour

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		dnsTicker := time.NewTicker(dnsInterval)
		defer dnsTicker.Stop()
		for {
			select {
			case <-ticker.C:
				_, _ = s.service.Benchmark(TriggerScheduler)
			case <-dnsTicker.C:
				_, _ = s.service.BenchmarkDNS(TriggerScheduler)
			case <-stop:
				return
			}
		}
	}()
}

func (s *BenchmarkScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopLocked()
}

func (s *BenchmarkScheduler) stopLocked() {
	if s.stop == nil {
		return
	}
	close(s.stop)
	s.stop = nil
}
