package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const maxTargets = 5

var errInvalidSettings = errors.New("settings must contain five named connectivity targets with valid urls")

func DefaultSettings() Settings {
	return Settings{
		IntervalMinutes: 10,
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

	if !strings.HasPrefix(strings.ToLower(raw), "http://") &&
		!strings.HasPrefix(strings.ToLower(raw), "https://") {
		raw = "https://" + raw
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() == "" || !strings.Contains(parsed.Hostname(), ".") {
		return "", false
	}

	if parsed.Path == "" || parsed.Path == "/" {
		if parsed.RawQuery == "" {
			return parsed.Scheme + "://" + parsed.Host + "/favicon.ico", true
		}
		parsed.Path = "/"
	}

	return parsed.String(), true
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
	if len(settings.Targets) != maxTargets {
		return Settings{}, errInvalidSettings
	}

	targets := make([]Target, 0, maxTargets)
	for index, target := range settings.Targets {
		name := strings.TrimSpace(target.Name)
		normalizedURL, ok := NormalizeTargetURL(target.URL)
		if name == "" || !ok {
			return Settings{}, errInvalidSettings
		}

		parsed, _ := url.Parse(normalizedURL)
		host := strings.TrimPrefix(parsed.Hostname(), "www.")
		id := strings.TrimSpace(target.ID)
		if id == "" {
			id = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		}
		if id == "" {
			id = "target-" + string(rune('1'+index))
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
