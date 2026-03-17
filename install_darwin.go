//go:build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const launchAgentLabel = "com.claude-usage-tray-go"

func installDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin")
}

func binaryName() string {
	return "claude-usage-tray-go"
}

func launchAgentDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents")
}

func launchAgentPath() string {
	return filepath.Join(launchAgentDir(), launchAgentLabel+".plist")
}

func setupAutostart(binPath string) error {
	dir := launchAgentDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>StandardOutPath</key>
    <string>/tmp/claude-usage-tray-go.stdout.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/claude-usage-tray-go.stderr.log</string>
    <key>ProcessType</key>
    <string>Interactive</string>
</dict>
</plist>
`, launchAgentLabel, binPath)

	if err := os.WriteFile(launchAgentPath(), []byte(plist), 0644); err != nil {
		return err
	}

	// Load the agent immediately
	exec.Command("launchctl", "load", launchAgentPath()).Run()
	return nil
}

func removeAutostart() error {
	plistPath := launchAgentPath()
	exec.Command("launchctl", "unload", plistPath).Run()

	err := os.Remove(plistPath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func autostartConfigured() bool {
	_, err := os.Stat(launchAgentPath())
	return err == nil
}
