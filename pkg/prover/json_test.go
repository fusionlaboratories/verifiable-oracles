package prover_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/prover"
)

var _jsonProver = prover.JsonProver{}

func Test_JsonProver_ProveVerify(t *testing.T) {
	for name, transcript := range transcriptTable {
		t.Run(name, func(t *testing.T) {
			var (
				assert     = assert.New(t)
				proof, err = _jsonProver.Prove(transcript)
			)

			assert.Nil(err)
			assert.NotNil(proof)

			result, err := _jsonProver.Verify(transcript, proof)
			assert.Nil(err)
			assert.True(result)
		})
	}
}

func Test_JsonProver_Verify_NilProof(t *testing.T) {
	var (
		assert = assert.New(t)

		transcript prover.Transcript
		proof      prover.Proof
	)

	result, err := _jsonProver.Verify(transcript, proof)

	assert.False(result)
	assert.NotNil(err)
	assert.NotEmpty(err.Error())
}
