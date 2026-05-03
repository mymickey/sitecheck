package main

import "testing"

func TestParseMyIPTrace(t *testing.T) {
	report, err := parseMyIPTrace("ip=13.231.213.207\nloc=JP\n")
	if err != nil {
		t.Fatalf("parseMyIPTrace() error = %v", err)
	}
	if report.IP != "13.231.213.207" {
		t.Fatalf("parseMyIPTrace() ip = %q, want %q", report.IP, "13.231.213.207")
	}
	if report.CountryCode != "JP" {
		t.Fatalf("parseMyIPTrace() countryCode = %q, want %q", report.CountryCode, "JP")
	}
}

func TestParseMyIPTrace_MissingFields(t *testing.T) {
	if _, err := parseMyIPTrace("fl=1009f26\nh=1.0.0.1\n"); err == nil {
		t.Fatal("parseMyIPTrace() error = nil, want error")
	}
}
