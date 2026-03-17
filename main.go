package main

import (
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {}

func onReady() {
	systray.SetTitle("claude ⚡...")
	systray.SetTooltip("Claude Code Stats " + Version)

	item5h := systray.AddMenuItem("5h: loading...", "5-hour rate limit")
	item7d := systray.AddMenuItem("7d: loading...", "7-day rate limit")
	itemExtra := systray.AddMenuItem("", "Extra usage")
	itemExtra.Hide()

	item5h.Disable()
	item7d.Disable()
	itemExtra.Disable()

	systray.AddSeparator()
	mRefresh := systray.AddMenuItem("↻ Refresh", "Refresh stats")
	mQuit := systray.AddMenuItem("✕ Quit", "Quit")

	updateUI := func(s StatsData) {
		systray.SetIcon(GenerateIconPNG(s.FiveHourPct))
		systray.SetTitle(FormatTitle(s))
		items := FormatMenuItems(s)
		if len(items) > 0 {
			item5h.SetTitle(items[0])
		}
		if len(items) > 1 {
			item7d.SetTitle(items[1])
		}
		if len(items) > 2 {
			itemExtra.SetTitle(items[2])
			itemExtra.Show()
		} else {
			itemExtra.Hide()
		}
	}

	// Initial load
	go func() {
		s := FetchStats(false)
		updateUI(s)
	}()

	// 60s auto-refresh ticker
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateUI(FetchStats(false))
		}
	}()

	// Click handler
	go func() {
		for {
			select {
			case <-mRefresh.ClickedCh:
				updateUI(FetchStats(true))
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
