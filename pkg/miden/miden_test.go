package miden_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/miden"
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

func TestMidenVersion(t *testing.T) {
	if !testHasMiden {
		t.Skip("miden not found, skipping")
	}

	assert := assert.New(t)

	v, err := miden.Version()
	assert.Nil(err)
	assert.Equal("Miden 0.6.0", v)
}
