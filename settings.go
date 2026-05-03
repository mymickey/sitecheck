package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const minTargets = 5

var errInvalidSettings = errors.New("settings must contain at least five connectivity targets with valid urls")
var errDuplicateTargetURL = errors.New("settings must not contain duplicate connectivity target urls")

func DefaultSettings() Settings {
	return Settings{
		IntervalMinutes:  10,
		DNSIntervalHours: 1,
		Targets: []Target{
			{
				ID:      "wechat",
				Name:    "WeChat",
				URL:     "https://res.wx.qq.com/a/wx_fed/assets/res/NTI4MWU5.ico",
				IconURL: "https://favicon.im/weixin.qq.com",
			},
			{
				ID:      "taobao",
				Name:    "Taobao",
				URL:     "https://www.taobao.com/favicon.ico",
				IconURL: "https://favicon.im/taobao.com",
			},
			{
				ID:      "google",
				Name:    "Google",
				URL:     "https://www.google.com/favicon.ico",
				IconURL: "https://favicon.im/google.com",
			},
			{
				ID:      "cloudflare",
				Name:    "Cloudflare",
				URL:     "https://www.cloudflare.com/favicon.ico",
				IconURL: "https://favicon.im/cloudflare.com",
			},
			{
				ID:      "github",
				Name:    "GitHub",
				URL:     "https://github.com/favicon.ico",
				IconURL: "https://favicon.im/github.com",
			},
		},
	}
}

func DefaultSettingsPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "sitecheck-settings.json"
	}
	return filepath.Join(configDir, "SiteCheck", "settings.json")
}

func NormalizeTargetURL(input string) (string, bool) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return "", false
	}

	addedScheme := false
	if !strings.HasPrefix(strings.ToLower(raw), "http://") &&
		!strings.HasPrefix(strings.ToLower(raw), "https://") {
		raw = "https://" + raw
		addedScheme = true
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() == "" || !strings.Contains(parsed.Hostname(), ".") {
		return "", false
	}

	if addedScheme && (parsed.Path == "" || parsed.Path == "/") {
		if parsed.RawQuery == "" {
			return parsed.Scheme + "://" + parsed.Host + "/favicon.ico", true
		}
		parsed.Path = "/"
		return parsed.String(), true
	}

	return strings.TrimSpace(input), true
}

type SettingsStore struct {
	path string
	mu   sync.Mutex
}

func NewSettingsStore(path string) *SettingsStore {
	return &SettingsStore{path: path}
}

func (s *SettingsStore) Load() (Settings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultSettings(), nil
	}
	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Settings{}, err
	}

	settings, err = normalizeSettings(settings)
	if err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func (s *SettingsStore) Save(settings Settings) (Settings, error) {
	normalized, err := normalizeSettings(settings)
	if err != nil {
		return Settings{}, err
	}

	data, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return Settings{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return Settings{}, err
	}
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return Settings{}, err
	}
	return normalized, nil
}

func normalizeSettings(settings Settings) (Settings, error) {
	if settings.IntervalMinutes <= 0 {
		settings.IntervalMinutes = 10
	}
	if settings.DNSIntervalHours < 1 || settings.DNSIntervalHours > 99 {
		settings.DNSIntervalHours = 1
	}
	if len(settings.Targets) < minTargets {
		return Settings{}, errInvalidSettings
	}

	targets := make([]Target, 0, len(settings.Targets))
	seenIDs := make(map[string]int, len(settings.Targets))
	seenURLs := make(map[string]struct{}, len(settings.Targets))
	for index, target := range settings.Targets {
		rawURL := strings.TrimSpace(target.URL)
		if _, exists := seenURLs[rawURL]; exists {
			return Settings{}, errDuplicateTargetURL
		}
		seenURLs[rawURL] = struct{}{}

		normalizedURL, ok := NormalizeTargetURL(rawURL)
		if !ok {
			return Settings{}, errInvalidSettings
		}

		parsed, _ := url.Parse(normalizedURL)
		host := strings.TrimPrefix(parsed.Hostname(), "www.")
		name := strings.TrimSpace(target.Name)
		if name == "" {
			name = host
		}
		if name == "" {
			return Settings{}, errInvalidSettings
		}
		id := strings.TrimSpace(target.ID)
		if id == "" {
			id = normalizedTargetID(name)
		}
		id = uniqueTargetID(id, seenIDs)
		if id == "" {
			id = uniqueTargetID("target-"+strconv.Itoa(index+1), seenIDs)
		}

		iconURL := strings.TrimSpace(target.IconURL)
		if iconURL == "" {
			iconURL = "https://favicon.im/" + host
		}

		targets = append(targets, Target{
			ID:      id,
			Name:    name,
			URL:     normalizedURL,
			IconURL: iconURL,
		})
	}

	settings.Targets = targets
	return settings, nil
}

func normalizedTargetID(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "-")
	var out strings.Builder
	lastDash := false
	for _, char := range name {
		switch {
		case char >= 'a' && char <= 'z', char >= '0' && char <= '9':
			out.WriteRune(char)
			lastDash = false
		case char == '.' || char == '-' || char == '_':
			if !lastDash && out.Len() > 0 {
				out.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(out.String(), "-")
}

func uniqueTargetID(base string, seen map[string]int) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "target"
	}
	count := seen[base]
	seen[base] = count + 1
	if count == 0 {
		return base
	}
	return base + "-" + strconv.Itoa(count+1)
}
