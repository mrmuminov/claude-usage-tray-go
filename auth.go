package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GetOAuthToken tries 4 sources: env → keychain (macOS) → credentials file → secret-tool (Linux)
func GetOAuthToken() (string, bool) {
	// 1. Environment variable
	if token := os.Getenv("CLAUDE_CODE_OAUTH_TOKEN"); token != "" {
		return token, true
	}

	// 2. Platform-specific keychain
	if token, ok := getFromKeychain(); ok {
		return token, true
	}

	// 3. ~/.claude/.credentials.json
	if token, ok := getFromCredentialsFile(); ok {
		return token, true
	}

	// 4. Platform-specific secret store
	if token, ok := getFromSecretStore(); ok {
		return token, true
	}

	return "", false
}

func getFromCredentialsFile() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	path := filepath.Join(home, ".claude", ".credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return parseCredentialsBlob(data)
}

func parseCredentialsBlob(data []byte) (string, bool) {
	var blob struct {
		ClaudeAiOauth struct {
			AccessToken string `json:"accessToken"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal(data, &blob); err != nil {
		return "", false
	}
	token := blob.ClaudeAiOauth.AccessToken
	if token == "" {
		return "", false
	}
	return token, true
}
