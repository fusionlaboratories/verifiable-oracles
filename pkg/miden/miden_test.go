package miden_test

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"github.com/qredo/verifiable-oracles/pkg/miden"
)

var testHasMiden bool

func out(v ...field.Element) field.Vector {
	l := len(v)
	if l < 16 {
		pad := make(field.Vector, 16-l)
		v = append(v, pad...)
	}
	return v
}

var midenTable = map[string]struct {
	assembly  []string
	inputFile miden.InputFile
	expected  field.Vector
}{
	"empty program": {
		assembly: []string{
			"begin",
			"end",
		},
		expected: out(),
	},
	"assert": {
		assembly: []string{
			"begin",
			"assert",
			"end",
		},
		inputFile: miden.InputFile{
			OperandStack: field.Vector{field.One()},
		},
		expected: out(),
	},
	"assertz": {
		assembly: []string{
			"begin",
			"assertz",
			"end",
		},
		inputFile: miden.InputFile{
			OperandStack: field.Vector{{}},
		},
		expected: out(),
	},
}

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

	for name, tc := range midenTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assembly := strings.Join(tc.assembly, "\n")
			output, _, err := miden.Run(assembly, tc.inputFile)

			handleExitError(t, err)
			assert.Equal(tc.expected, output)
		})
	}

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
