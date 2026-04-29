package main

const (
	StatusWait        = "Waiting"
	StatusAvailable   = "Available"
	StatusUnavailable = "Unavailable"
)

type Target struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"iconUrl"`
}

type Settings struct {
	IntervalMinutes int      `json:"intervalMinutes"`
	Targets         []Target `json:"targets"`
}

type ProbeResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	IconURL   string `json:"iconUrl"`
	Status    string `json:"status"`
	LatencyMS int    `json:"latencyMs"`
	Error     string `json:"error,omitempty"`
	CheckedAt string `json:"checkedAt"`
}

type BenchmarkSummary struct {
	FastestMS  int    `json:"fastestMs"`
	SlowestMS  int    `json:"slowestMs"`
	HasResults bool   `json:"hasResults"`
	CheckedAt  string `json:"checkedAt"`
}

type BenchmarkReport struct {
	Results []ProbeResult    `json:"results"`
	Summary BenchmarkSummary `json:"summary"`
}
