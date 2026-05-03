package main

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateIPAPIPrefix(t *testing.T) {
	prefix := generateIPAPIPrefix()
	if len(prefix) != 32 {
		t.Fatalf("generateIPAPIPrefix() length = %d, want 32", len(prefix))
	}
	if prefix[13:23] != "jason5ng32" {
		t.Fatalf("generateIPAPIPrefix() fixed segment = %q, want %q", prefix[13:23], "jason5ng32")
	}
	for _, r := range prefix[:13] {
		if r < '0' || r > '9' {
			t.Fatalf("generateIPAPIPrefix() timestamp segment = %q, want milliseconds digits", prefix[:13])
		}
	}
	for _, r := range prefix[23:] {
		if !(r >= '0' && r <= '9') && !(r >= 'a' && r <= 'z') {
			t.Fatalf("generateIPAPIPrefix() random segment = %q, want lowercase base36", prefix[23:])
		}
	}
}

func TestGenerateSurfsharkPrefix(t *testing.T) {
	prefix := generateSurfsharkPrefix()
	if !strings.HasPrefix(prefix, "jn32") {
		t.Errorf("generateSurfsharkPrefix() = %q, want prefix 'jn32'", prefix)
	}
	if len(prefix) != 13 {
		t.Fatalf("generateSurfsharkPrefix() length = %d, want 13", len(prefix))
	}
}

func TestNormalizeIPAPI(t *testing.T) {
	tests := []struct {
		name    string
		data    IPAPIResponse
		want    DNSCheckpoint
		wantErr bool
	}{
		{
			name: "edns priority",
			data: IPAPIResponse{
				DNS:  &IPAPIEntry{Geo: "China - Telecom", IP: "1.1.1.1"},
				EDNS: &IPAPIEntry{Geo: "US - Comcast", IP: "2.2.2.2"},
			},
			want: DNSCheckpoint{
				Result: DNSResult{Country: "US", ISP: "Comcast", IP: "2.2.2.2"},
			},
			wantErr: false,
		},
		{
			name: "fallback to dns",
			data: IPAPIResponse{
				DNS:  &IPAPIEntry{Geo: "China - Telecom", IP: "1.1.1.1"},
				EDNS: nil,
			},
			want: DNSCheckpoint{
				Result: DNSResult{Country: "China", ISP: "Telecom", IP: "1.1.1.1"},
			},
			wantErr: false,
		},
		{
			name: "no separation format",
			data: IPAPIResponse{
				DNS: &IPAPIEntry{Geo: "Localhost", IP: "127.0.0.1"},
			},
			want: DNSCheckpoint{
				Result: DNSResult{Country: "Unknown", ISP: "Localhost", IP: "127.0.0.1"},
			},
			wantErr: false,
		},
		{
			name: "empty response",
			data: IPAPIResponse{},
			want: DNSCheckpoint{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeIPAPI(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeIPAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("normalizeIPAPI() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestRunDNSTest_Timeout(t *testing.T) {
	// A quick unit test that ensures a timeout results in "Timeout" status
	// We'll pass a context that is already canceled to simulate a timeout immediately.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	report := RunDNSTest(ctx)
	if len(report.Checkpoints) != 4 {
		t.Fatalf("RunDNSTest() shape = %d checkpoints, want exactly 4", len(report.Checkpoints))
	}

	for _, cp := range report.Checkpoints {
		if cp.Result.ISP != "Timeout" && cp.Result.ISP != "Failed" && cp.Result.ISP != "--" {
			t.Errorf("Expected Timeout, Failed, or --, got %s", cp.Result.ISP)
		}
	}
}
