package miden

import (
	"os/exec"
	"strings"
)

// Execute Miden and get it's version
func Version() (string, error) {
	out, err := exec.Command("miden", "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
