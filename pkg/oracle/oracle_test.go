package oracle

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// Oracle Input Table
var table = []struct {
	blockNumber uint64
	hash        common.Hash
}{
	{blockNumber: uint64(0), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 0})},
	{blockNumber: uint64(1), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 1})},
	{blockNumber: uint64(2), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 2})},
	{blockNumber: uint64(3), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 3})},
	{blockNumber: uint64(4), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 4})},
	{blockNumber: uint64(5), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 5})},
	{blockNumber: uint64(6), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 6})},
	{blockNumber: uint64(7), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 7})},
	{blockNumber: uint64(8), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 8})},
	{blockNumber: uint64(9), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 9})},
	{blockNumber: uint64(10), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 10})},
	{blockNumber: uint64(11), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 11})},
	{blockNumber: uint64(12), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 12})},
	{blockNumber: uint64(13), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 13})},
	{blockNumber: uint64(14), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 14})},
	{blockNumber: uint64(15), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 15})},
}

func TestAddBlock(t *testing.T) {
	o := NewInMemoryOracle()

	for _, tc := range table {

		o.AddBlock(tc.blockNumber, tc.hash)
	}
}

func TestGetBlockHash(t *testing.T) {
	assert := assert.New(t)

	o := NewInMemoryOracle()

	for _, tc := range table {

		o.AddBlock(tc.blockNumber, tc.hash)
	}

	for _, tc := range table {
		h, err := o.GetBlockHash(tc.blockNumber)

		assert.Nil(err)
		assert.Equal(tc.hash, h)
	}
}

func TestMissingBlockError(t *testing.T) {
	assert := assert.New(t)

	o := NewInMemoryOracle()
	number := uint64(132)

	_, err := o.GetBlockHash(number)

	var missingBlockError *MissingBlockError
	assert.ErrorAs(err, &missingBlockError)
	assert.Equal(number, missingBlockError.BlockNumber)
}

// benchmark result.  Added to prevent compiler form optimizing away the
// benchmarks.
var _result common.Hash

func BenchmarkInMemory(t *testing.B) {
	o := NewInMemoryOracle()

	for n := 0; n < t.N; n++ {
		for _, tc := range table {
			o.AddBlock(tc.blockNumber, tc.hash)
		}
	}

	_result, _ = o.GetBlockHash(uint64(0))
}

// Transcript which is subset of table
var transcriptTable = []struct {
	blockNumber uint64
	hash        common.Hash
}{
	table[13],
	table[5],
	table[1],
	table[2],
	table[7],
	table[3],
	table[11],
}

func TestTranscript(t *testing.T) {
	assert := assert.New(t)

	o := NewInMemoryOracle()
	for _, tc := range table {
		o.AddBlock(tc.blockNumber, tc.hash)
	}

	to := NewTranscriptOracle(o)

	for _, tc := range transcriptTable {
		h, err := to.GetBlockHash(tc.blockNumber)
		assert.Equal(tc.hash, h)
		assert.Nil(err)
	}

	transcript := to.GetTranscript()
	assert.Equal(len(transcriptTable), len(transcript))

	for i := range transcript {
		assert.Equal(transcriptTable[i].blockNumber, transcript[i].BlockNumber)
		assert.Equal(transcriptTable[i].hash, transcript[i].BlockHash)
	}
}

func BenchmarkTranscript(t *testing.B) {
	o := NewInMemoryOracle()

	for _, tc := range table {
		o.AddBlock(tc.blockNumber, tc.hash)
	}

	to := NewTranscriptOracle(o)

	r := common.Hash{}

	for n := 0; n < t.N; n++ {
		for _, tc := range transcriptTable {
			r, _ = to.GetBlockHash(tc.blockNumber)
		}

		transcript := to.GetTranscript()

		if len(transcript) > 0 {
			r = transcript[0].BlockHash
		}
	}

	_result = r
}
