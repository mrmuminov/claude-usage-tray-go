package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	apiURL      = "https://api.anthropic.com/api/oauth/usage"
	cacheMaxAge = 60 * time.Second
)

var cacheFile = filepath.Join(os.TempDir(), "claude", "tray-cache.json")

type StatsData struct {
	FiveHourPct   int
	FiveHourReset string
	SevenDayPct   int
	SevenDayReset string
	ExtraEnabled  bool
	ExtraPct      int
	ExtraUsed     float64
	ExtraLimit    float64
}

type usageWindow struct {
	Utilization *float64 `json:"utilization"`
	ResetsAt    *string  `json:"resets_at"`
}

type extraUsage struct {
	IsEnabled    *bool    `json:"is_enabled"`
	Utilization  *float64 `json:"utilization"`
	UsedCredits  *float64 `json:"used_credits"`
	MonthlyLimit *float64 `json:"monthly_limit"`
}

type usageResponse struct {
	FiveHour   *usageWindow `json:"five_hour"`
	SevenDay   *usageWindow `json:"seven_day"`
	ExtraUsage *extraUsage  `json:"extra_usage"`
}

// FetchStats fetches usage stats with 60s cache. Set force=true to bypass cache.
func FetchStats(force bool) StatsData {
	// Try fresh cache first (only if not forced)
	if !force {
		if cached := loadCache(); cached != nil {
			log.Println("data: using fresh cache")
			return parseUsage(*cached)
		}
	}

	// Try API
	if usage := fetchFromAPI(); usage != nil {
		log.Println("data: fetched from API")
		return parseUsage(*usage)
	}

	// Fall back to stale cache only on auto-refresh, not on manual force
	if !force {
		if stale := loadCacheAny(); stale != nil {
			log.Println("data: using stale cache fallback")
			return parseUsage(*stale)
		}
	} else {
		log.Println("data: force refresh — API failed, returning empty")
	}

	return StatsData{}
}

func loadCache() *usageResponse {
	info, err := os.Stat(cacheFile)
	if err != nil {
		return nil
	}
	if time.Since(info.ModTime()) >= cacheMaxAge {
		return nil
	}
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}
	var u usageResponse
	if err := json.Unmarshal(data, &u); err != nil {
		return nil
	}
	return &u
}

func loadCacheAny() *usageResponse {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}
	var u usageResponse
	if err := json.Unmarshal(data, &u); err != nil {
		return nil
	}
	return &u
}

func fetchFromAPI() *usageResponse {
	token, ok := GetOAuthToken()
	if !ok {
		log.Println("data: no OAuth token found")
		return nil
	}
	log.Printf("data: token found (len=%d)", len(token))

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("data: build request error: %v", err)
		return nil
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	req.Header.Set("User-Agent", "claude-code/2.1.34")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("data: HTTP error: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("data: API returned %d: %s", resp.StatusCode, truncate(string(body), 200))
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("data: read body error: %v", err)
		return nil
	}

	// Validate expected field exists
	var check map[string]json.RawMessage
	if err := json.Unmarshal(body, &check); err != nil {
		log.Printf("data: JSON parse error: %v", err)
		return nil
	}
	if _, ok := check["five_hour"]; !ok {
		log.Printf("data: unexpected response — no 'five_hour' field: %s", truncate(string(body), 200))
		return nil
	}

	// Save to cache
	_ = os.MkdirAll(filepath.Dir(cacheFile), 0755)
	_ = os.WriteFile(cacheFile, body, 0644)

	var u usageResponse
	if err := json.Unmarshal(body, &u); err != nil {
		log.Printf("data: parse usageResponse error: %v", err)
		return nil
	}
	return &u
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return fmt.Sprintf("%s...[%d bytes]", s[:n], len(s))
}

func parseUsage(u usageResponse) StatsData {
	var s StatsData

	if u.FiveHour != nil {
		if u.FiveHour.Utilization != nil {
			s.FiveHourPct = int(*u.FiveHour.Utilization + 0.5)
		}
		if u.FiveHour.ResetsAt != nil {
			s.FiveHourReset = *u.FiveHour.ResetsAt
		}
	}

	if u.SevenDay != nil {
		if u.SevenDay.Utilization != nil {
			s.SevenDayPct = int(*u.SevenDay.Utilization + 0.5)
		}
		if u.SevenDay.ResetsAt != nil {
			s.SevenDayReset = *u.SevenDay.ResetsAt
		}
	}

	if u.ExtraUsage != nil {
		if u.ExtraUsage.IsEnabled != nil {
			s.ExtraEnabled = *u.ExtraUsage.IsEnabled
		}
		if u.ExtraUsage.Utilization != nil {
			s.ExtraPct = int(*u.ExtraUsage.Utilization + 0.5)
		}
		if u.ExtraUsage.UsedCredits != nil {
			s.ExtraUsed = *u.ExtraUsage.UsedCredits / 100.0
		}
		if u.ExtraUsage.MonthlyLimit != nil {
			s.ExtraLimit = *u.ExtraUsage.MonthlyLimit / 100.0
		}
	}

	return s
}
