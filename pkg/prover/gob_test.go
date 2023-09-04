package prover_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/prover"
)

func Test_GobProver_Prove_Empty(t *testing.T) {
	assert := assert.New(t)
	transcript := prover.Transcript{}
	g := &prover.GobProver{}
	proof, err := g.Prove(transcript)

	assert.Nil(err)
	assert.NotNil(proof)
}

func Test_GobProver_ProveVerify(t *testing.T) {
	assert := assert.New(t)
	transcript := prover.Transcript{}
	g := &prover.GobProver{}

	proof, err := g.Prove(transcript)
	assert.Nil(err)

	result, err := g.Verify(transcript, proof)
	assert.Nil(err)
	assert.True(result)
}

func Test_GobProver_Prove_MutateTranscript(t *testing.T) {
	assert := assert.New(t)
	transcript := prover.Transcript{}
	g := &prover.GobProver{}

	proof, err := g.Prove(transcript)
	assert.Nil(err)

	transcript2 := prover.Transcript{
		BlockNumber: uint64(1),
	}

	result, err := g.Verify(transcript2, proof)

	assert.Nil(err)
	assert.False(result)
}

func Test_GobProver_Verify_NilProof(t *testing.T) {
	assert := assert.New(t)
	transcript := prover.Transcript{}
	g := &prover.GobProver{}
	proof := prover.Proof{}

	result, err := g.Verify(transcript, proof)

	assert.False(result)
	assert.NotNil(err)
	assert.NotEmpty(err.Error())
}
