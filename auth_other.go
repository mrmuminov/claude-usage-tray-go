//go:build !darwin && !linux

package main

func getFromKeychain() (string, bool) {
	return "", false
}

func getFromSecretStore() (string, bool) {
	return "", false
}
