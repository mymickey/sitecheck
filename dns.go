package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// IPAPIResponse represents the response from ip-api edns endpoint
type IPAPIResponse struct {
	DNS  *IPAPIEntry `json:"dns"`
	EDNS *IPAPIEntry `json:"edns"`
}

type IPAPIEntry struct {
	Geo string `json:"geo"`
	IP  string `json:"ip"`
}

// SurfsharkEntry represents the response from Surfshark DNS endpoint
type SurfsharkEntry struct {
	ISP         string `json:"ISP"`
	Country     string `json:"Country"`
	City        string `json:"City"`
	IP          string `json:"IP"`
	Leak        bool   `json:"Leak"`
	CountryCode string `json:"CountryCode"`
}

func generateRandomBase36() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000000000))
	return strconv.FormatInt(n.Int64(), 36)
}

func isTimeoutError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

func generateIPAPIPrefix() string {
	// current unix time string + jason5ng32 + random base36 substring
	now := strconv.FormatInt(time.Now().Unix(), 10)
	random := generateRandomBase36()
	return fmt.Sprintf("%sjason5ng32%s", now, random)
}

func generateSurfsharkPrefix() string {
	// jn32 + random base36 substring
	random := generateRandomBase36()
	return fmt.Sprintf("jn32%s", random)
}

func fetchIPAPI(ctx context.Context, client *http.Client) (DNSCheckpoint, error) {
	prefix := generateIPAPIPrefix()
	urlStr := fmt.Sprintf("https://%s.edns.ip-api.com/json?lang=en", prefix)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return DNSCheckpoint{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return DNSCheckpoint{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DNSCheckpoint{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var data IPAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return DNSCheckpoint{}, err
	}

	return normalizeIPAPI(data)
}

func normalizeIPAPI(data IPAPIResponse) (DNSCheckpoint, error) {
	var entry *IPAPIEntry
	if data.EDNS != nil && data.EDNS.Geo != "" {
		entry = data.EDNS
	} else if data.DNS != nil && data.DNS.Geo != "" {
		entry = data.DNS
	} else {
		return DNSCheckpoint{}, fmt.Errorf("no valid dns or edns data")
	}

	country := "Unknown"
	isp := entry.Geo

	parts := strings.SplitN(entry.Geo, " - ", 2)
	if len(parts) == 2 {
		country = strings.TrimSpace(parts[0])
		isp = strings.TrimSpace(parts[1])
	}

	return DNSCheckpoint{
		Result: DNSResult{
			ISP:     isp,
			IP:      entry.IP,
			Country: country,
		},
	}, nil
}

func fetchSurfshark(ctx context.Context, client *http.Client) ([]DNSCheckpoint, error) {
	prefix := generateSurfsharkPrefix()
	urlStr := fmt.Sprintf("https://%s.ipv4.surfsharkdns.com/", prefix)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// We need to parse into a map, but we must preserve original JSON object ordering
	// Go's json.Unmarshal into map[string]interface{} doesn't preserve order.
	// But the rules state: "preserve original traversal order".
	// Let's decode manually using json.RawMessage or use a slice if we can.
	// Surfshark returns an object where keys are IPs. We can decode the keys sequentially using json.Decoder.
	dec := json.NewDecoder(resp.Body)
	t, err := dec.Token()
	if err != nil || t != json.Delim('{') {
		return nil, fmt.Errorf("expected JSON object")
	}

	var checkpoints []DNSCheckpoint
	for dec.More() {
		_, err := dec.Token() // key (IP address)
		if err != nil {
			return nil, err
		}

		var entry SurfsharkEntry
		if err := dec.Decode(&entry); err != nil {
			return nil, err
		}

		if len(checkpoints) < 3 {
			checkpoints = append(checkpoints, DNSCheckpoint{
				Result: DNSResult{
					ISP:     entry.ISP,
					IP:      entry.IP,
					Country: entry.Country,
				},
			})
		}
	}

	return checkpoints, nil
}

// RunDNSTest executes both DNS checks concurrently and returns the combined report
func RunDNSTest(ctx context.Context) DNSTestReport {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Ensure overall timeout
	testCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	type ipapiResult struct {
		cp  DNSCheckpoint
		err error
	}
	type ssResult struct {
		cps []DNSCheckpoint
		err error
	}

	ipapiChan := make(chan ipapiResult, 1)
	ssChan := make(chan ssResult, 1)

	go func() {
		cp, err := fetchIPAPI(testCtx, client)
		ipapiChan <- ipapiResult{cp, err}
	}()

	go func() {
		cps, err := fetchSurfshark(testCtx, client)
		ssChan <- ssResult{cps, err}
	}()

	var finalCheckpoints = make([]DNSCheckpoint, 4)
	for i := 0; i < 4; i++ {
		finalCheckpoints[i] = DNSCheckpoint{
			Name: fmt.Sprintf("check point #%d", i+1),
			Result: DNSResult{ISP: "--", IP: "--", Country: "--"},
		}
	}

	ipapiRes := <-ipapiChan
	if ipapiRes.err != nil {
		status := "Failed"
		if isTimeoutError(ipapiRes.err) {
			status = "Timeout"
		}
		finalCheckpoints[0].Result = DNSResult{ISP: status, IP: status, Country: status}
	} else {
		finalCheckpoints[0].Result = ipapiRes.cp.Result
	}

	ssRes := <-ssChan
	if ssRes.err != nil {
		status := "Failed"
		if isTimeoutError(ssRes.err) {
			status = "Timeout"
		}
		for i := 1; i < 4; i++ {
			finalCheckpoints[i].Result = DNSResult{ISP: status, IP: status, Country: status}
		}
	} else {
		for i := 0; i < 3; i++ {
			if i < len(ssRes.cps) {
				finalCheckpoints[i+1].Result = ssRes.cps[i].Result
			}
		}
	}

	return DNSTestReport{
		Checkpoints: finalCheckpoints,
		CheckedAt:   time.Now().Format(time.RFC3339),
	}
}
