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

var (
	_testHasMiden  bool
	_defaultOutput = out()
)

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
	},
	"get field element from advice stack": {
		assembly: []string{
			"begin",
			"adv_push.1",
			"assert_eq",
			"end",
		},
		inputFile: miden.InputFile{
			OperandStack: field.Vector{field.One()},
			AdviceStack:  field.Vector{field.One()},
		},
	},
}

func init() {
	if _, err := exec.LookPath("miden"); err == nil {
		_testHasMiden = true
	}
}

func needsMiden(t *testing.T) {
	t.Helper()

	if !_testHasMiden {
		t.Skip("miden not found, skipping")
	}
}

func handleExitError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			t.Errorf(string(exitError.Stderr))
		} else {
			t.Errorf("unknown error %v", err)
		}

		return false
	}

	return true
}

func TestMiden(t *testing.T) {
	needsMiden(t)

	for name, tc := range midenTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assembly := strings.Join(tc.assembly, "\n")
			output, _, err := miden.Run(assembly, tc.inputFile)

			// Avoid cluttering test output by only checking output when
			// the execution was successful
			if handleExitError(t, err) {
				expectedOutput := tc.expected
				if expectedOutput == nil {
					expectedOutput = _defaultOutput
				}
				assert.Equal(expectedOutput, output)
			}
		})
	}

}

func TestMidenVersion(t *testing.T) {
	needsMiden(t)

	assert := assert.New(t)

	v, err := miden.Version()
	assert.Nil(err)
	assert.Equal("Miden 0.6.0", v)
}

func TestMidenRunFile(t *testing.T) {
	needsMiden(t)

	assert := assert.New(t)
	output, _, err := miden.RunFile("testdata/test.masm", "testdata/test.json")

	handleExitError(t, err)
	assert.Equal(make(field.Vector, 16), output)
}
