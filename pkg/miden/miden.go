package miden

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

// Execute Miden and get it's version
func Version() (string, error) {
	out, err := exec.Command("miden", "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func extractOuptut(outLines []string) (field.Vector, error) {
	var (
		outputPrefix = "Output: ["
		outputSuffix = "]"
		elemSep      = ", "

		output string
	)

	for _, line := range outLines {
		if strings.HasPrefix(line, outputPrefix) {
			output = line
			break
		}
	}

	if len(output) == 0 {
		return nil, errors.New("miden: output line not found")
	}

	output = strings.TrimPrefix(output, outputPrefix)
	output = strings.TrimSuffix(output, outputSuffix)

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

// TODO: Figure how to use it
func extractHash(outLines []string) ([]byte, error) {
	outputPrefix := "Executing program with hash "
	outputSuffix := "... done"

	var output string
	for _, line := range outLines {
		if strings.HasPrefix(line, outputPrefix) {
			output = line
			break
		}
	}

	if len(output) == 0 {
		return nil, errors.New("miden: hash line not found")
	}

	output = strings.TrimPrefix(output, outputPrefix)

	if suffIndex := strings.Index(output, outputSuffix); suffIndex != -1 {
		output = output[:suffIndex]
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

	output, err1 := extractOuptut(outLines)
	hash, err2 := extractHash(outLines)

	return output, hash, errors.Join(err1, err2)
}
