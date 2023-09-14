// Basic Miden driver for Golang
package miden

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
)

// TODO:
// - [ ] Use --output flag from miden,
// - [ ] Consider splitting the functionality into separate files,

// Treating ProgramHash as []byte for now
type ProgramHash = []byte

// Treating Proof as []byte for now
type Proof = []byte

// Execute Miden and get it's version
func Version(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "miden", "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func extractLine(lines []string, prefix string, suffix string) (string, bool) {
	outputIndex := slices.IndexFunc[[]string](lines,
		func(line string) bool { return strings.HasPrefix(line, prefix) })

	if outputIndex == -1 {
		return "", false
	}
	line := lines[outputIndex]

	// Trim prefix
	line = line[len(prefix):]

	// Trim after last occurrence of suffix
	if suffIndex := strings.LastIndex(line, suffix); suffIndex != -1 {
		line = line[:suffIndex]
	}

	return line, true
}

func extractHashRun(outLines []string) ([]byte, error) {
	outputPrefix := "Executing program with hash "
	outputSuffix := "... done"

	output, ok := extractLine(outLines, outputPrefix, outputSuffix)
	if !ok {
		return nil, errors.New("miden: hash line not found")
	}

	return hex.DecodeString(output)
}

func extractHashCompile(outLines []string) ([]byte, error) {
	outputPrefix := "program hash is "

	output, ok := extractLine(outLines, outputPrefix, "")
	if !ok {
		return nil, errors.New("miden: hash line not found")
	}

	return hex.DecodeString(output)
}

func tempFile(contents []byte, pattern string) (name string, err error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return
	}
	name = f.Name()
	if _, err = f.Write(contents); err != nil {
		return
	}
	if err = f.Close(); err != nil {
		return
	}
	return
}

func Run(ctx context.Context, assembly string, input Input) (output Output, hash ProgramHash, err error) {
	// Create a temporary directory for all the files instead
	dirPath, err := os.MkdirTemp("", "miden*")
	if err != nil {
		return
	}
	// cleanup
	defer os.RemoveAll(dirPath)

	var (
		assemblyPath = path.Join(dirPath, "assembly.masm")
		inputPath    = path.Join(dirPath, "input.json")
		outputPath   = path.Join(dirPath, "output.json")
	)

	// Writing assembly file
	err = os.WriteFile(assemblyPath, []byte(assembly), 0644)
	if err != nil {
		return
	}

	// Writing input file
	inputData, err := json.Marshal(input)
	if err != nil {
		return
	}
	err = os.WriteFile(inputPath, inputData, 0644)
	if err != nil {
		return
	}

	// running file
	hash, err = RunFile(ctx, assemblyPath, inputPath, outputPath)
	if err != nil {
		return
	}

	// getting output
	outputBytes, err := os.ReadFile(outputPath)
	if err != nil {
		return
	}
	err = json.Unmarshal(outputBytes, &output)
	return
}

func RunFile(ctx context.Context, assemblyPath string, inputPath string, outputPath string) (ProgramHash, error) {
	cmd := exec.CommandContext(ctx, "miden", "run", "--assembly", assemblyPath, "--input", inputPath, "--output", outputPath)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(string(out), "\n")

	return extractHashRun(outLines)
}

func Compile(ctx context.Context, assembly string) (ProgramHash, error) {
	assemblyFile, err := tempFile([]byte(assembly), "*.masm")
	if err != nil {
		return nil, err
	}
	defer os.Remove(assemblyFile)

	return CompileFile(ctx, assemblyFile)
}

func CompileFile(ctx context.Context, assemblyPath string) (ProgramHash, error) {
	cmd := exec.CommandContext(ctx, "miden", "compile", "--assembly", assemblyPath)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(string(out), "\n")

	hash, err := extractHashCompile(outLines)
	return hash, err
}

func Prove(ctx context.Context, assembly string, input Input) (Proof, error) {
	panic("unimplemented")

}

func ProveFile(ctx context.Context, assmeblyPath string, inputPath string, proofPath string, outputPath string) error {
	cmd := exec.CommandContext(ctx, "miden", "prove", "--assembly", assmeblyPath, "--input", inputPath, "--proof", proofPath, "--output", outputPath)

	_, err := cmd.Output()

	return err
}

func Verify(ctx context.Context, programHash ProgramHash, proof Proof, input Input) (bool, error) {
	panic("unimplemented")
}

// TODO: Does this need output as well?
func VerifyFile(ctx context.Context, programHash ProgramHash, proofPath string, inputPath string) (bool, error) {
	panic("unimplemented")
}
