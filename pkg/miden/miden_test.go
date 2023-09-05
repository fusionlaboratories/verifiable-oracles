package miden_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
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

func TestMidenRun(t *testing.T) {
	assert := assert.New(t)
	output, _, err := miden.Run("testdata/test.masm")

	assert.Nil(err)
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			t.Log(string(exitError.Stderr))
		} else {
			t.Logf("unknown error %v", err)
		}
	}
	assert.Equal(make(field.Vector, 16), output)
}
