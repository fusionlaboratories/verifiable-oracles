package prover_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/prover"
)

type DummyProver struct {
}

type Input struct{}

func (*DummyProver) Prove(input Input) (prover.Proof, error) {
	return prover.Proof{}, nil
}

func (*DummyProver) Verify(input Input, proof prover.Proof) (bool, error) {
	return true, nil
}

var _ prover.Prover[Input] = (*DummyProver)(nil)
var _ prover.Verifier[Input] = (*DummyProver)(nil)

func Test_DummyProver_Prove(t *testing.T) {
	var (
		assert = assert.New(t)
		dummy  = &DummyProver{}
	)

	_, err := dummy.Prove(Input{})
	assert.Nil(err)
}

func Test_DummyProver_Verify(t *testing.T) {
	var (
		assert = assert.New(t)
		dummy  = &DummyProver{}
	)

	_, err := dummy.Verify(Input{}, prover.Proof{})
	assert.Nil(err)
}
