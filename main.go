package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

//go:embed logo.png
var logoPNG []byte

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if isInstalled() {
				fmt.Print("Already installed. Reinstall? [y/N]: ")
				reader := bufio.NewReader(os.Stdin)
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer != "y" && answer != "yes" {
					fmt.Println("Cancelled.")
					return
				}
				fmt.Println("Reinstalling...")
			}
			if err := Install(); err != nil {
				fmt.Fprintf(os.Stderr, "Install failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Installed successfully.")
			return
		case "uninstall":
			if err := Uninstall(); err != nil {
				fmt.Fprintf(os.Stderr, "Uninstall failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Uninstalled successfully.")
			return
		case "status":
			PrintStatus()
			return
		case "version", "--version", "-v":
			fmt.Printf("claude-usage-tray-go %s\n", Version)
			return
		case "help", "--help", "-h":
			printHelp()
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
			fmt.Fprintln(os.Stderr)
			printHelp()
			os.Exit(1)
		}
	}

	systray.Run(onReady, onExit)
}

func printHelp() {
	fmt.Printf("Claude Usage Tray %s\n", Version)
	fmt.Println()
	fmt.Println("Usage: claude-usage-tray-go [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  (none)       Start the tray application")
	fmt.Println("  install      Install binary and configure autostart")
	fmt.Println("  uninstall    Remove autostart and installed binary")
	fmt.Println("  status       Show installation status")
	fmt.Println("  version      Print version")
	fmt.Println("  help         Show this help message")
}

func onExit() {}

func onReady() {
	systray.SetIcon(ResizeLogoPNG(logoPNG, 64))
	systray.SetTitle("claude ⚡...")
	systray.SetTooltip("Claude Usage Tray " + Version)

	item5h := systray.AddMenuItem("5h: loading...", "5-hour rate limit")
	item7d := systray.AddMenuItem("7d: loading...", "7-day rate limit")
	itemExtra := systray.AddMenuItem("", "Extra usage")
	itemExtra.Hide()

	item5h.Disable()
	item7d.Disable()
	itemExtra.Disable()

	systray.AddSeparator()
	mRefresh := systray.AddMenuItem("Refresh", "Refresh stats")
	mRefresh.SetIcon(GenerateMenuActionIcon("refresh"))
	mGitHub := systray.AddMenuItem("GitHub", "https://github.com/mrmuminov/claude-usage-tray-go")
	mGitHub.SetIcon(GenerateMenuActionIcon("github"))
	mQuit := systray.AddMenuItem("Quit", "Quit")
	mQuit.SetIcon(GenerateMenuActionIcon("quit"))

	updateUI := func(s StatsData) {
		systray.SetTooltip(fmt.Sprintf("Claude Usage Tray — 5h: %d%% | 7d: %d%%", s.FiveHourPct, s.SevenDayPct))
		systray.SetTitle(FormatTitle(s))
		items := FormatMenuItems(s)
		if len(items) > 0 {
			item5h.SetTitle(items[0])
			item5h.SetIcon(GenerateMenuDotIcon(s.FiveHourPct))
		}
		if len(items) > 1 {
			item7d.SetTitle(items[1])
			item7d.SetIcon(GenerateMenuDotIcon(s.SevenDayPct))
		}
		if len(items) > 2 {
			itemExtra.SetTitle(items[2])
			itemExtra.SetIcon(GenerateMenuDotIcon(s.ExtraPct))
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
			case <-mGitHub.ClickedCh:
				openBrowser("https://github.com/mrmuminov/claude-usage-tray-go")
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
