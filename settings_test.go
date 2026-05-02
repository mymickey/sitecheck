package main

import (
	"errors"
	"testing"
)

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

func TestNormalizeSettingsAllowsCustomTargetsBeyondDefaultFive(t *testing.T) {
	settings := DefaultSettings()
	settings.Targets = append(settings.Targets, Target{
		Name: "example.com",
		URL:  "https://example.com/status",
	})

	got, err := normalizeSettings(settings)
	if err != nil {
		t.Fatalf("normalizeSettings() error = %v, want nil", err)
	}

	if len(got.Targets) != 6 {
		t.Fatalf("len(normalizeSettings().Targets) = %d, want 6", len(got.Targets))
	}

	last := got.Targets[5]
	if last.ID == "" || last.Name == "" || last.URL == "" || last.IconURL == "" {
		t.Fatalf("custom target = %+v, want populated fields", last)
	}
	if last.Name != "example.com" {
		t.Fatalf("custom target name = %q, want %q", last.Name, "example.com")
	}
	if last.IconURL != "https://favicon.im/example.com" {
		t.Fatalf("custom target icon = %q, want %q", last.IconURL, "https://favicon.im/example.com")
	}
}

func TestNormalizeSettingsRejectsDuplicateTargetURLs(t *testing.T) {
	settings := DefaultSettings()
	settings.Targets = append(settings.Targets, Target{
		Name: "GitHub Duplicate",
		URL:  "https://github.com/favicon.ico",
	})

	_, err := normalizeSettings(settings)
	if !errors.Is(err, errDuplicateTargetURL) {
		t.Fatalf("normalizeSettings() error = %v, want %v", err, errDuplicateTargetURL)
	}
}

func TestNormalizeSettingsAllowsTrimDistinctURLs(t *testing.T) {
	settings := DefaultSettings()
	settings.Targets = append(settings.Targets, Target{
		Name: "GitHub Root",
		URL:  "https://github.com",
	})

	got, err := normalizeSettings(settings)
	if err != nil {
		t.Fatalf("normalizeSettings() error = %v, want nil", err)
	}

	last := got.Targets[len(got.Targets)-1]
	if last.URL != "https://github.com" {
		t.Fatalf("custom target url = %q, want %q", last.URL, "https://github.com")
	}
}
