package miden

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"golang.org/x/exp/slices"
)

// Execute Miden and get it's version
func Version() (string, error) {
	out, err := exec.Command("miden", "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func extractLine(lines []string, prefix string, suffix string) (string, bool) {
	outputIndex := slices.IndexFunc[string](lines,
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

func extractOutput(outLines []string) (field.Vector, error) {
	var (
		outputPrefix = "Output: ["
		outputSuffix = "]"
		elemSep      = ", "
	)

	output, ok := extractLine(outLines, outputPrefix, outputSuffix)
	if !ok {
		return nil, errors.New("miden: output line not found")
	}

	outElems := strings.Split(output, elemSep)
	result := make(field.Vector, len(outElems))

	for i := 0; i < len(outElems); i++ {
		eStr := outElems[i]
		if _, err := result[i].SetString(eStr); err != nil {
			return nil, err
		}
	}

	return result, nil
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

func Run(assembly string, input InputFile) (field.Vector, []byte, error) {
	assemblyFile, err := tempFile([]byte(assembly), "*.masm")
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(assemblyFile)

	inputContents, err := json.Marshal(input)
	if err != nil {
		return nil, nil, err

	}

	inputFile, err := tempFile(inputContents, "*.json")
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(inputFile)

	return RunFile(assemblyFile, inputFile)
}

func RunFile(assemblyPath string, inputPath string) (field.Vector, []byte, error) {
	cmd := exec.Command("miden", "run", "--assembly", assemblyPath, "--input", inputPath)

	out, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}

	outLines := strings.Split(string(out), "\n")

	output, err1 := extractOutput(outLines)
	hash, err2 := extractHashRun(outLines)

	return output, hash, errors.Join(err1, err2)
}

func Compile(assembly string) ([]byte, error) {
	assemblyFile, err := tempFile([]byte(assembly), "*.masm")
	if err != nil {
		return nil, err
	}
	defer os.Remove(assemblyFile)

	return CompileFile(assemblyFile)
}

func CompileFile(assemblyPath string) ([]byte, error) {
	cmd := exec.Command("miden", "compile", "--assembly", assemblyPath)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(string(out), "\n")

	hash, err := extractHashCompile(outLines)
	return hash, err
}
