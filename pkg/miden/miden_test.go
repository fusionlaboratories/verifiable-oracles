package miden_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
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
			hash, output, err := miden.Run(context.Background(), assembly, tc.inputFile)

			// Avoid cluttering test output by only checking output when
			// the execution was successful
			if handleExitError(t, err) {
				expectedOutput := tc.expected
				if expectedOutput == nil {
					expectedOutput = _defaultOutput
				}
				assert.Equal(tc.hash, hashToHex(hash))
				assert.Equal(expectedOutput, output.Stack)
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

func TestMidenProve(t *testing.T) {
	needsMiden(t)

	for name, tc := range midenTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assembly := strings.Join(tc.assembly, "\n")

			hash, output, proof, err := miden.Prove(context.Background(), assembly, tc.inputFile)

			if handleExitError(t, err) {
				expectedOutput := tc.expected
				if expectedOutput == nil {
					expectedOutput = _defaultOutput
				}
				assert.Equal(tc.hash, hashToHex(hash))
				assert.Equal(expectedOutput, output.Stack)
				assert.NotEmpty(proof)
			}
		})
	}
}

func TestMidenVerify(t *testing.T) {
	needsMiden(t)

	for name, tc := range midenTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assembly := strings.Join(tc.assembly, "\n")

			hash, output, proof, err := miden.Prove(context.Background(), assembly, tc.inputFile)

			if handleExitError(t, err) {
				result, err := miden.Verify(context.Background(), hash, proof, tc.inputFile, output)
				if handleExitError(t, err) {
					assert.True(result)
				}
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

	var (
		assert = assert.New(t)

		assemblyPath = "testdata/test.masm"
		inputPath    = "testdata/input.json"
		outputPath   = "testdata/run/output.json"

		expectedHash = "a4820838f4914083b432faaaef596a86b84c6a061d0bf90711d6ba294244e308"
	)

	hash, err := miden.RunFile(context.Background(), assemblyPath, inputPath, outputPath)
	handleExitError(t, err)
	assert.Equal(expectedHash, hashToHex(hash))

	data, err := os.ReadFile(outputPath)
	assert.Nil(err)

	var output miden.Output
	err = json.Unmarshal(data, &output)
	assert.Nil(err)
	assert.Equal(out(), output.Stack)
}

func TestMidenCompileFile(t *testing.T) {
	needsMiden(t)

	var (
		assert = assert.New(t)

		assemblyPath = "testdata/test.masm"

		expectedHash = "a4820838f4914083b432faaaef596a86b84c6a061d0bf90711d6ba294244e308"
	)

	hash, err := miden.CompileFile(context.Background(), assemblyPath)

	handleExitError(t, err)
	assert.Equal(expectedHash, hashToHex(hash))
}

func TestMidenProveFile(t *testing.T) {
	needsMiden(t)

	var (
		assemblyPath = "testdata/test.masm"
		inputPath    = "testdata/input.json"
		outputPath   = "testdata/prove/output.json"
		proofPath    = "testdata/prove/proof.bin"
	)

	_, err := miden.ProveFile(context.Background(), assemblyPath, inputPath, outputPath, proofPath)
	handleExitError(t, err)
}

func TestMidenVerifyFile(t *testing.T) {
	needsMiden(t)

	var (
		assemblyPath = "testdata/test.masm"
		inputPath    = "testdata/input.json"
		outputPath   = "testdata/verify/output.json"
		proofPath    = "testdata/verify/proof.bin"

		expectedHash = "a4820838f4914083b432faaaef596a86b84c6a061d0bf90711d6ba294244e308"
	)

	programHash, _ := hex.DecodeString(expectedHash)

	_, err := miden.ProveFile(context.Background(), assemblyPath, inputPath, outputPath, proofPath)
	if handleExitError(t, err) {

		result, err := miden.VerifyFile(context.Background(), programHash, inputPath, outputPath, proofPath)
		if handleExitError(t, err) {
			assert.True(t, result)
		}
	}
}
