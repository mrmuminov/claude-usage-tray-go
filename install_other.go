//go:build !linux && !darwin && !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func installDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin")
}

func binaryName() string {
	return "claude-usage-tray-go"
}

func setupAutostart(binPath string) error {
	return fmt.Errorf("autostart not supported on %s", runtime.GOOS)
}

func removeAutostart() error {
	return fmt.Errorf("autostart not supported on %s", runtime.GOOS)
}

func autostartConfigured() bool {
	return false
}
