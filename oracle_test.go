package oracle

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// Oracle Input Table
var table = []struct {
	blockNumber *big.Int
	hash        common.Hash
}{
	{blockNumber: big.NewInt(0), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 0})},
	{blockNumber: big.NewInt(1), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 1})},
	{blockNumber: big.NewInt(2), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 2})},
	{blockNumber: big.NewInt(3), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 3})},
	{blockNumber: big.NewInt(4), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 4})},
	{blockNumber: big.NewInt(5), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 5})},
	{blockNumber: big.NewInt(6), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 6})},
	{blockNumber: big.NewInt(7), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 7})},
	{blockNumber: big.NewInt(8), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 8})},
	{blockNumber: big.NewInt(9), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 9})},
	{blockNumber: big.NewInt(10), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 10})},
	{blockNumber: big.NewInt(11), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 11})},
	{blockNumber: big.NewInt(12), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 12})},
	{blockNumber: big.NewInt(13), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 13})},
	{blockNumber: big.NewInt(14), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 14})},
	{blockNumber: big.NewInt(15), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 15})},
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
	number := big.NewInt(123)

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

	_result, _ = o.GetBlockHash(big.NewInt(0))
}

// Transcript which is subset of table
var transcriptTable = []struct {
	blockNumber *big.Int
	hash        common.Hash
}{
	{blockNumber: big.NewInt(13), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 13})},
	{blockNumber: big.NewInt(5), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 5})},
	{blockNumber: big.NewInt(1), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 1})},
	{blockNumber: big.NewInt(2), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 2})},
	{blockNumber: big.NewInt(7), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 7})},
	{blockNumber: big.NewInt(3), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 3})},
	{blockNumber: big.NewInt(11), hash: common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 11})},
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
		assert.Equal(transcriptTable[i].blockNumber, transcript[i].Number)
		assert.Equal(transcriptTable[i].hash, transcript[i].Hash)
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
			r = transcript[0].Hash
		}
	}

	_result = r
}
