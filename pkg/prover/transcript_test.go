package prover_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/prover"
)

// A simple transcript table
var transcriptTable = map[string]prover.Transcript{
	"empty":            {},
	"blockNumber":      {BlockNumber: uint64(1)},
	"blockHash":        {BlockHash: common.BytesToHash([]byte{1})},
	"transactionIndex": {TransactionIndex: uint64(1)},
	"transactionHash":  {TransactionHash: common.BytesToHash([]byte{1})},
}

func TestTranscript_Marshal(t *testing.T) {
	for name, transcript := range transcriptTable {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(transcript)

			assert.Nil(t, err)
			assert.NotEmpty(t, b)
		})
	}
}

func TestTranscript_Marshal_Roundtrip(t *testing.T) {
	for name, transcript := range transcriptTable {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(transcript)
			assert.Nil(t, err)

			var transcript2 prover.Transcript
			err = json.Unmarshal(b, &transcript2)

			assert.Nil(t, err)
			assert.Equal(t, transcript, transcript2)
		})
	}
}
