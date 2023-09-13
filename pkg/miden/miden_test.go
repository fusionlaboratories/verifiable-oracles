package miden_test

import (
	"context"
	"encoding/hex"
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
	inputFile miden.Input
	expected  field.Vector
	hash      string
}{
	"empty program": {
		assembly: []string{
			"begin",
			"end",
		},
		hash: "f0db3924f3e2d677a51924b09ecef8a12416a6ceb09fadd39785bb4f685cab66",
	},
	"assert": {
		assembly: []string{
			"begin",
			"assert",
			"end",
		},
		inputFile: miden.Input{
			OperandStack: field.Vector{field.One()},
		},
		hash: "1858ec2e6abdf1d1447474e5ab8e1313c4f93276e82f3baac9a056d6ecdc0c9b",
	},
	"assertz": {
		assembly: []string{
			"begin",
			"assertz",
			"end",
		},
		inputFile: miden.Input{
			OperandStack: field.Vector{{}},
		},
		hash: "f9b9df59a9549b8e8833d86ec3f1f97f0fdfe002c24eb8661c4d7242d3c14a45",
	},
	"add one to two": {
		assembly: []string{
			"begin",
			"add",
			"end",
		},
		inputFile: miden.Input{
			OperandStack: field.Vector{field.One(), field.NewElement(2)},
		},
		expected: out(field.NewElement(3)),
		hash:     "63c2b2b5cf6abd6414fb93cc7af4ad22fed1c8d3182ea1a01d3aba005c453c57",
	},
	"get field element from advice stack": {
		assembly: []string{
			"begin",
			"adv_push.1",
			"assert_eq",
			"end",
		},
		inputFile: miden.Input{
			OperandStack: field.Vector{field.One()},
			AdviceStack:  field.Vector{field.One()},
		},
		hash: "5857c99e44517e8b7bf8abc514041ab75590b425147be13eae69be9bda411db7",
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

func hashToHex(hash miden.ProgramHash) string {

	return hex.EncodeToString(hash)
}

func TestMidenRun(t *testing.T) {
	needsMiden(t)

	for name, tc := range midenTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assembly := strings.Join(tc.assembly, "\n")
			output, hash, err := miden.Run(context.Background(), assembly, tc.inputFile)

			// Avoid cluttering test output by only checking output when
			// the execution was successful
			if handleExitError(t, err) {
				expectedOutput := tc.expected
				if expectedOutput == nil {
					expectedOutput = _defaultOutput
				}
				assert.Equal(tc.hash, hashToHex(hash))
				assert.Equal(expectedOutput, output)
			}
		})
	}
}

func TestMidenCompile(t *testing.T) {
	needsMiden(t)

	for name, tc := range midenTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assembly := strings.Join(tc.assembly, "\n")
			hash, err := miden.Compile(context.Background(), assembly)

			// Avoid cluttering test output by only checking output when
			// the execution was successful
			if handleExitError(t, err) {
				assert.Equal(tc.hash, hashToHex(hash))
			}
		})
	}
}

func TestMidenVersion(t *testing.T) {
	needsMiden(t)

	assert := assert.New(t)

	v, err := miden.Version(context.Background())
	assert.Nil(err)
	assert.Equal("Miden 0.6.0", v)
}

func TestMidenRunFile(t *testing.T) {
	needsMiden(t)

	assert := assert.New(t)
	output, hash, err := miden.RunFile(context.Background(), "testdata/test.masm", "testdata/input.json")

	handleExitError(t, err)
	assert.Equal("a4820838f4914083b432faaaef596a86b84c6a061d0bf90711d6ba294244e308", hashToHex(hash))
	assert.Equal(make(field.Vector, 16), output)
}

func TestMidenCompileFile(t *testing.T) {
	needsMiden(t)

	assert := assert.New(t)
	hash, err := miden.CompileFile(context.Background(), "testdata/test.masm")

	handleExitError(t, err)
	assert.Equal("a4820838f4914083b432faaaef596a86b84c6a061d0bf90711d6ba294244e308", hashToHex(hash))
}
