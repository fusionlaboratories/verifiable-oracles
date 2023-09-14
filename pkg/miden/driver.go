package miden

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

// Miden driver
type driver struct {
	wd  string
	err error
}

func newTmpDirDriver() *driver {
	wd, err := os.MkdirTemp("", "miden*")
	return &driver{wd: wd, err: err}
}

// paths
func (d *driver) assemblyPath() string {
	return path.Join(d.wd, "assembly.masm")
}

func (d *driver) inputPath() string {
	return path.Join(d.wd, "input.json")
}

func (d *driver) outputPath() string {
	return path.Join(d.wd, "output.json")
}

func (d *driver) proofPath() string {
	return path.Join(d.wd, "proof.bin")
}

func (d *driver) setAssembly(assembly string) error {
	if d.err != nil {
		return d.err
	}
	d.err = os.WriteFile(d.assemblyPath(), []byte(assembly), 0644)
	return d.err
}

func (d *driver) setInput(input Input) error {
	if d.err != nil {
		return d.err
	}

	var data []byte
	data, d.err = json.Marshal(input)
	if d.err != nil {
		return d.err
	}

	d.err = os.WriteFile(d.inputPath(), data, 0644)
	return d.err
}

func (d *driver) output() (Output, error) {
	var (
		output Output
		data   []byte
	)

	if d.err != nil {
		return output, d.err
	}

	data, d.err = os.ReadFile(d.outputPath())
	if d.err != nil {
		return output, d.err
	}
	d.err = json.Unmarshal(data, &output)
	return output, d.err
}

func (d *driver) setOutput(output Output) error {
	if d.err != nil {
		return d.err
	}

	var data []byte
	data, d.err = json.Marshal(output)
	if d.err != nil {
		return d.err
	}

	d.err = os.WriteFile(d.outputPath(), data, 0644)
	return d.err
}

func (d *driver) proof() (Proof, error) {
	var proof Proof

	if d.err != nil {
		return proof, d.err
	}

	proof, d.err = os.ReadFile(d.proofPath())
	return proof, d.err
}

func (d *driver) setProof(proof Proof) error {
	if d.err != nil {
		return d.err
	}

	d.err = os.WriteFile(d.proofPath(), proof, 0644)
	return d.err
}

func (d *driver) cleanup() error {
	if d.wd == "" {
		return errors.New("trying to remove an empty wd")
	}
	return os.RemoveAll(d.wd)
}
