// Basic Miden driver for Golang
package miden

import (
	"context"
	"encoding/hex"
	"errors"
	"os/exec"
	"slices"
	"strings"
)

// TODO:
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

func extractHashProve(outLines []string) ([]byte, error) {
	outputPrefix := "Program with hash "
	outputSuffix := " proved in"

	output, ok := extractLine(outLines, outputPrefix, outputSuffix)
	if !ok {
		return nil, errors.New("miden: hash line not found")
	}

	return hex.DecodeString(output)
}

// Trying to reimplement Run with driver
func Run(ctx context.Context, assembly string, input Input) (hash ProgramHash, output Output, err error) {
	d := newTmpDirDriver()
	defer d.cleanup()

	if err = d.setAssembly(assembly); err != nil {
		return
	}
	if err = d.setInput(input); err != nil {
		return
	}

	if hash, err = RunFile(ctx, d.assemblyPath(), d.inputPath(), d.outputPath()); err != nil {
		return
	}

	// getting output
	output, err = d.output()
	return
}

func Compile(ctx context.Context, assembly string) (ProgramHash, error) {
	d := newTmpDirDriver()
	defer d.cleanup()

	if err := d.setAssembly(assembly); err != nil {
		return nil, err
	}

	return CompileFile(ctx, d.assemblyPath())
}

func Prove(ctx context.Context, assembly string, input Input) (hash ProgramHash, output Output, proof Proof, err error) {
	d := newTmpDirDriver()
	defer d.cleanup()

	if err = d.setAssembly(assembly); err != nil {
		return
	}
	if err = d.setInput(input); err != nil {
		return
	}

	if hash, err = ProveFile(ctx, d.assemblyPath(), d.inputPath(), d.outputPath(), d.proofPath()); err != nil {
		return
	}
	if output, err = d.output(); err != nil {
		return
	}

	proof, err = d.proof()
	return
}

func Verify(ctx context.Context, programHash ProgramHash, proof Proof, input Input, output Output) (r bool, err error) {
	d := newTmpDirDriver()
	defer d.cleanup()

	if err = d.setInput(input); err != nil {
		return
	}
	if err = d.setOutput(output); err != nil {
		return
	}
	if err = d.setProof(proof); err != nil {
		return
	}

	return VerifyFile(ctx, programHash, d.inputPath(), d.outputPath(), d.proofPath())
}

func CompileFile(ctx context.Context, assemblyPath string) (ProgramHash, error) {
	cmd := exec.CommandContext(ctx, "miden", "compile", "--assembly", assemblyPath)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(string(out), "\n")
	return extractHashCompile(outLines)
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

func ProveFile(ctx context.Context, assemblyPath string, inputPath string, outputPath string, proofPath string) (ProgramHash, error) {
	cmd := exec.CommandContext(ctx, "miden", "prove", "--assembly", assemblyPath, "--input", inputPath, "--output", outputPath, "--proof", proofPath)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(string(out), "\n")

	return extractHashProve(outLines)
}

func VerifyFile(ctx context.Context, programHash ProgramHash, inputPath string, outputPath string, proofPath string) (bool, error) {
	hash := hex.EncodeToString(programHash)
	cmd := exec.CommandContext(ctx, "miden", "verify", "--program-hash", hash, "--input", inputPath, "--output", outputPath, "--proof", proofPath)
	_, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return true, nil
}
