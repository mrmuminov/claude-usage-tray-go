//go:build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const registryKeyPath = `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
const registryValueName = "ClaudeUsageTray"

func installDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		home, _ := os.UserHomeDir()
		localAppData = filepath.Join(home, "AppData", "Local")
	}
	return filepath.Join(localAppData, "claude-usage-tray-go")
}

func binaryName() string {
	return "claude-usage-tray-go.exe"
}

func setupAutostart(binPath string) error {
	return exec.Command("reg", "add", registryKeyPath,
		"/v", registryValueName,
		"/t", "REG_SZ",
		"/d", binPath,
		"/f",
	).Run()
}

func removeAutostart() error {
	err := exec.Command("reg", "delete", registryKeyPath,
		"/v", registryValueName,
		"/f",
	).Run()
	if err != nil {
		out, queryErr := exec.Command("reg", "query", registryKeyPath,
			"/v", registryValueName,
		).CombinedOutput()
		if queryErr != nil && strings.Contains(string(out), "ERROR") {
			return nil
		}
		return err
	}
	return nil
}

func autostartConfigured() bool {
	err := exec.Command("reg", "query", registryKeyPath,
		"/v", registryValueName,
	).Run()
	return err == nil
}
