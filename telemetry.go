package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	aptabaseAppKey = "A-EU-7884746208"
	appVersion     = "0.1.0"
)

type TelemetryClient struct {
	appKey     string
	host       string
	httpClient *http.Client
	sessionID  string
	osVersion  string
	locale     string
}

type telemetryEvent struct {
	Timestamp   string         `json:"timestamp"`
	SessionID   string         `json:"sessionId"`
	EventName   string         `json:"eventName"`
	SystemProps telemetryProps `json:"systemProps"`
	Props       telemetryProps `json:"props,omitempty"`
}

type telemetryProps map[string]any

func NewTelemetryClient(appKey string) *TelemetryClient {
	if strings.TrimSpace(appKey) == "" {
		return nil
	}

	return &TelemetryClient{
		appKey:     appKey,
		host:       telemetryHost(appKey),
		httpClient: &http.Client{Timeout: 4 * time.Second},
		sessionID:  newTelemetrySessionID(),
		osVersion:  macOSVersion(),
		locale:     telemetryLocale(),
	}
}

func (t *TelemetryClient) Track(eventName string, props telemetryProps) {
	if t == nil || eventName == "" {
		return
	}

	payload := []telemetryEvent{{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		SessionID: t.sessionID,
		EventName: eventName,
		SystemProps: telemetryProps{
			"locale":     t.locale,
			"osName":     telemetryOSName(),
			"osVersion":  t.osVersion,
			"isDebug":    false,
			"appVersion": appVersion,
			"sdkVersion": "sitecheck-go@0.1.0",
		},
		Props: props,
	}}

	go t.send(payload)
}

func (t *TelemetryClient) send(events []telemetryEvent) {
	body, err := json.Marshal(events)
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodPost, t.host+"/api/v0/events", bytes.NewReader(body))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("App-Key", t.appKey)

	response, err := t.httpClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
}

func telemetryHost(appKey string) string {
	if strings.Contains(strings.ToUpper(appKey), "EU") {
		return "https://eu.aptabase.com"
	}
	return "https://us.aptabase.com"
}

func newTelemetrySessionID() string {
	return time.Now().UTC().Format("20060102150405") + randomDigits(8)
}

func randomDigits(length int) string {
	const digits = "0123456789"
	buffer := make([]byte, length)
	if _, err := crand.Read(buffer); err != nil {
		return strings.Repeat("0", length)
	}
	var out strings.Builder
	out.Grow(length)
	for _, value := range buffer {
		out.WriteByte(digits[int(value)%len(digits)])
	}
	return out.String()
}

func telemetryLocale() string {
	locale := strings.TrimSpace(os.Getenv("LANG"))
	if locale == "" {
		return "en-US"
	}
	locale = strings.Split(locale, ".")[0]
	locale = strings.ReplaceAll(locale, "_", "-")
	return locale
}

func telemetryOSName() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}

func macOSVersion() string {
	if runtime.GOOS != "darwin" {
		return ""
	}
	output, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
