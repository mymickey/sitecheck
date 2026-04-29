package main

import "testing"

func TestDefaultSettingsUsesTenMinuteIntervalAndFiveTargets(t *testing.T) {
	settings := DefaultSettings()

	if settings.IntervalMinutes != 10 {
		t.Errorf("DefaultSettings().IntervalMinutes = %d, want 10", settings.IntervalMinutes)
	}
	if len(settings.Targets) != 5 {
		t.Fatalf("len(DefaultSettings().Targets) = %d, want 5", len(settings.Targets))
	}

	for _, target := range settings.Targets {
		if target.ID == "" || target.Name == "" || target.URL == "" || target.IconURL == "" {
			t.Errorf("DefaultSettings() target = %+v, want populated ID, Name, URL, and IconURL", target)
		}
	}
}

func TestNormalizeTargetURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		ok    bool
	}{
		{
			name:  "bare domain uses https favicon",
			input: "github.com",
			want:  "https://github.com/favicon.ico",
			ok:    true,
		},
		{
			name:  "existing path is preserved",
			input: "https://example.com/api/health",
			want:  "https://example.com/api/health",
			ok:    true,
		},
		{
			name:  "query path is preserved",
			input: "https://example.com/?v=1",
			want:  "https://example.com/?v=1",
			ok:    true,
		},
		{
			name:  "hostname must look routable",
			input: "localhost",
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := NormalizeTargetURL(tt.input)
			if ok != tt.ok {
				t.Fatalf("NormalizeTargetURL(%q) ok = %t, want %t", tt.input, ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("NormalizeTargetURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
