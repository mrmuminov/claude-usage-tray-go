package main

import (
	"fmt"
	"strings"
	"time"
)

// FormatTitle returns "claude ⚡N%" using the 5-hour rate.
func FormatTitle(s StatsData) string {
	return fmt.Sprintf("claude ⚡%d%%", s.FiveHourPct)
}

// FormatMenuItems returns the display lines for menu items.
func FormatMenuItems(s StatsData) []string {
	var items []string

	items = append(items, fmt.Sprintf("5h Rate:  %s %d%%  %s",
		buildBar(s.FiveHourPct, 8),
		s.FiveHourPct,
		formatResetTime(s.FiveHourReset, "time"),
	))

	items = append(items, fmt.Sprintf("7d Rate:  %s %d%%  %s",
		buildBar(s.SevenDayPct, 8),
		s.SevenDayPct,
		formatResetTime(s.SevenDayReset, "date"),
	))

	if s.ExtraEnabled {
		items = append(items, fmt.Sprintf("Extra:    %s $%.2f/$%.2f",
			buildBar(s.ExtraPct, 8),
			s.ExtraUsed,
			s.ExtraLimit,
		))
	}

	return items
}

func buildBar(pct, width int) string {
	if pct > 100 {
		pct = 100
	}
	if pct < 0 {
		pct = 0
	}
	filled := pct * width / 100
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

// formatResetTime parses an ISO 8601 timestamp and formats it.
// style "time" → "3:45pm", style "date" → "Jan 2"
func formatResetTime(iso, style string) string {
	if iso == "" || iso == "null" {
		return ""
	}

	t, err := parseISO(iso)
	if err != nil {
		return ""
	}

	switch style {
	case "time":
		// Format like "3:45pm" — use local time
		return strings.ToLower(t.Local().Format("3:04pm"))
	case "date":
		// Format like "Jan 2"
		return t.Local().Format("Jan 2")
	default:
		return t.Local().Format("Jan 2 3:04pm")
	}
}

// parseISO tries multiple ISO 8601 formats.
func parseISO(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05.999999999Z07:00",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as ISO 8601", s)
}
