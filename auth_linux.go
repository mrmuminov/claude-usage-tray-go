//go:build linux

package main

import (
	"os/exec"
	"strings"
)

func getFromKeychain() (string, bool) {
	return "", false
}

func getFromSecretStore() (string, bool) {
	out, err := exec.Command("timeout", "2", "secret-tool", "lookup", "service", "Claude Code-credentials").Output()
	if err != nil {
		return "", false
	}
	blob := strings.TrimSpace(string(out))
	return parseCredentialsBlob([]byte(blob))
}
