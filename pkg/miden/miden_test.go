package miden_test

import (
	"os/exec"
	"testing"
)

var testHasMiden bool

func init() {
	if _, err := exec.LookPath("miden"); err == nil {
		testHasMiden = true
	}
}

func TestMiden(t *testing.T) {
	if !testHasMiden {
		t.Skip("miden not found, skipping")
	}
}
