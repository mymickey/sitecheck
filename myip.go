package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const myIPURL = "https://1.0.0.1/cdn-cgi/trace"

func parseMyIPTrace(body string) (MyIPReport, error) {
	var (
		ip  string
		loc string
	)

	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "ip":
			ip = strings.TrimSpace(value)
		case "loc":
			loc = strings.ToUpper(strings.TrimSpace(value))
		}
	}

	if err := scanner.Err(); err != nil {
		return MyIPReport{}, err
	}

	if ip == "" && loc == "" {
		return MyIPReport{}, fmt.Errorf("trace response missing ip and loc")
	}

	return MyIPReport{
		IP:          ip,
		CountryCode: loc,
		CheckedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

func fetchMyIPOnce(ctx context.Context, client *http.Client) (MyIPReport, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, myIPURL, nil)
	if err != nil {
		return MyIPReport{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return MyIPReport{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MyIPReport{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var lines []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return MyIPReport{}, err
	}

	return parseMyIPTrace(strings.Join(lines, "\n"))
}

func FetchMyIPUntilSuccess(ctx context.Context) (MyIPReport, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	for {
		report, err := fetchMyIPOnce(ctx, client)
		if err == nil {
			return report, nil
		}
		if ctx.Err() != nil {
			return MyIPReport{}, ctx.Err()
		}
		timer := time.NewTimer(10 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return MyIPReport{}, ctx.Err()
		case <-timer.C:
		}
	}
}
