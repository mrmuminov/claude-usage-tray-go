//go:build darwin

package main

import (
	"os/exec"
	"strings"
)

func getFromKeychain() (string, bool) {
	out, err := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w").Output()
	if err != nil {
		return "", false
	}
	blob := strings.TrimSpace(string(out))
	return parseCredentialsBlob([]byte(blob))
}

func getFromSecretStore() (string, bool) {
	return "", false
}
