//go:build linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const appDesktopName = "claude-usage-tray-go.desktop"

func installDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin")
}

func binaryName() string {
	return "claude-usage-tray-go"
}

func autostartDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "autostart")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "autostart")
}

func setupAutostart(binPath string) error {
	dir := autostartDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	desktopEntry := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Claude Usage Tray
Comment=Claude Code usage stats in system tray
Exec=%s
Icon=utilities-system-monitor
Terminal=false
Categories=Utility;
StartupNotify=false
X-GNOME-Autostart-enabled=true
`, binPath)

	return os.WriteFile(filepath.Join(dir, appDesktopName), []byte(desktopEntry), 0644)
}

func removeAutostart() error {
	err := os.Remove(filepath.Join(autostartDir(), appDesktopName))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func autostartConfigured() bool {
	_, err := os.Stat(filepath.Join(autostartDir(), appDesktopName))
	return err == nil
}
