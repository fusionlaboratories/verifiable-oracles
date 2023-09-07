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

func handleExitError(t *testing.T, err error) {
	t.Helper()
	assert := assert.New(t)

	assert.Nil(err)
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			t.Log(string(exitError.Stderr))
		} else {
			t.Logf("unknown error %v", err)
		}
	}
}

func TestMiden(t *testing.T) {
	if !testHasMiden {
		t.Skip("miden not found, skipping")
	}
	assert := assert.New(t)

	assembly := `begin
end`

	output, _, err := miden.Run(assembly, miden.InputFile{})
	expectedOutput := make(field.Vector, 16)

	handleExitError(t, err)
	assert.Equal(expectedOutput, output)
}

func TestMiden_assert(t *testing.T) {
	if !testHasMiden {
		t.Skip("miden not found, skipping")
	}
	assert := assert.New(t)

	assembly := `begin
	assert
end`

	inputFile := miden.InputFile{
		OperandStack: field.Vector{field.One()},
	}

	output, _, err := miden.Run(assembly, inputFile)
	expectedOutput := make(field.Vector, 16)

	handleExitError(t, err)
	assert.Equal(expectedOutput, output)
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

func TestMidenRunFile(t *testing.T) {
	if !testHasMiden {
		t.Skip("miden not found, skipping")
	}

	assert := assert.New(t)
	output, _, err := miden.RunFile("testdata/test.masm", "testdata/test.json")

	handleExitError(t, err)
	assert.Equal(make(field.Vector, 16), output)
}
