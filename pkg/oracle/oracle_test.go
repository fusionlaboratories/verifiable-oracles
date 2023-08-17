package oracle_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/qredo/verifiable-oracles/pkg/oracle"
	"github.com/stretchr/testify/assert"
)

// Variables

// Sample blocks
var blockTable = []struct {
	blockNumber uint64
	blockHash   common.Hash
}{
	{blockNumber: uint64(0), blockHash: common.BytesToHash([]byte{0})},
	{blockNumber: uint64(1), blockHash: common.BytesToHash([]byte{1})},
	{blockNumber: uint64(2), blockHash: common.BytesToHash([]byte{2})},
	{blockNumber: uint64(3), blockHash: common.BytesToHash([]byte{3})},
	{blockNumber: uint64(4), blockHash: common.BytesToHash([]byte{4})},
	{blockNumber: uint64(5), blockHash: common.BytesToHash([]byte{5})},
	{blockNumber: uint64(6), blockHash: common.BytesToHash([]byte{6})},
	{blockNumber: uint64(7), blockHash: common.BytesToHash([]byte{7})},
	{blockNumber: uint64(8), blockHash: common.BytesToHash([]byte{8})},
	{blockNumber: uint64(9), blockHash: common.BytesToHash([]byte{9})},
	{blockNumber: uint64(10), blockHash: common.BytesToHash([]byte{10})},
	{blockNumber: uint64(11), blockHash: common.BytesToHash([]byte{11})},
	{blockNumber: uint64(12), blockHash: common.BytesToHash([]byte{12})},
	{blockNumber: uint64(13), blockHash: common.BytesToHash([]byte{13})},
	{blockNumber: uint64(14), blockHash: common.BytesToHash([]byte{14})},
	{blockNumber: uint64(15), blockHash: common.BytesToHash([]byte{15})},
}

// transcript table of getting block hashes
var blockTranscriptTable = []struct {
	blockNumber uint64
	blockHash   common.Hash
}{
	blockTable[13],
	blockTable[5],
	blockTable[1],
	blockTable[2],
	blockTable[7],
	blockTable[3],
	blockTable[11],
}

// Sample transactions
var transactionTable = []struct {
	blockHash        common.Hash
	transactionIndex uint64
	transactionHash  common.Hash
}{
	{blockHash: blockTable[0].blockHash, transactionIndex: uint64(0), transactionHash: common.BytesToHash([]byte{0})},
	{blockHash: blockTable[1].blockHash, transactionIndex: uint64(0), transactionHash: common.BytesToHash([]byte{1})},
	{blockHash: blockTable[3].blockHash, transactionIndex: uint64(0), transactionHash: common.BytesToHash([]byte{2})},
	{blockHash: blockTable[3].blockHash, transactionIndex: uint64(1), transactionHash: common.BytesToHash([]byte{3})},
	{blockHash: blockTable[5].blockHash, transactionIndex: uint64(0), transactionHash: common.BytesToHash([]byte{4})},
	{blockHash: blockTable[5].blockHash, transactionIndex: uint64(1), transactionHash: common.BytesToHash([]byte{5})},
	{blockHash: blockTable[5].blockHash, transactionIndex: uint64(2), transactionHash: common.BytesToHash([]byte{6})},
	{blockHash: blockTable[7].blockHash, transactionIndex: uint64(0), transactionHash: common.BytesToHash([]byte{7})},
	{blockHash: blockTable[7].blockHash, transactionIndex: uint64(1), transactionHash: common.BytesToHash([]byte{8})},
	{blockHash: blockTable[7].blockHash, transactionIndex: uint64(2), transactionHash: common.BytesToHash([]byte{9})},
	{blockHash: blockTable[7].blockHash, transactionIndex: uint64(3), transactionHash: common.BytesToHash([]byte{10})},
	{blockHash: blockTable[11].blockHash, transactionIndex: uint64(0), transactionHash: common.BytesToHash([]byte{11})},
	{blockHash: blockTable[11].blockHash, transactionIndex: uint64(1), transactionHash: common.BytesToHash([]byte{12})},
	{blockHash: blockTable[11].blockHash, transactionIndex: uint64(2), transactionHash: common.BytesToHash([]byte{13})},
	{blockHash: blockTable[11].blockHash, transactionIndex: uint64(3), transactionHash: common.BytesToHash([]byte{14})},
	{blockHash: blockTable[11].blockHash, transactionIndex: uint64(4), transactionHash: common.BytesToHash([]byte{15})},
}

// transcript table of getting transaction hashes
var transactionTranscriptTable = []struct {
	blockHash        common.Hash
	transactionIndex uint64
	transactionHash  common.Hash
}{
	transactionTable[7],
	transactionTable[2],
	transactionTable[13],
	transactionTable[3],
	transactionTable[11],
	transactionTable[5],
}

// benchmark result.  Added to prevent compiler form optimizing away the
// benchmarks.
var _result common.Hash

// Test adding Blocks to an Empty Oracle
func Test_InMemory_AddBlock(t *testing.T) {
	o := &oracle.InMemoryOracle{}

	for _, tc := range blockTable {

		o.AddBlock(tc.blockNumber, tc.blockHash)
	}
}

// Test if AddBlock makes a copy
func Test_InMemory_AddBlock_MakesCopy(t *testing.T) {
	o := &oracle.InMemoryOracle{}
	number := uint64(42)
	hash := common.BytesToHash([]byte{1})

	o.AddBlock(number, hash)
	hash.SetBytes([]byte{2})

	h, _ := o.GetBlockHash(number)

	assert.NotEqual(t, hash, h)
}

// Testing GetBlockHash
func Test_InMemory_GetBlockHash(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	for _, tc := range blockTable {

		o.AddBlock(tc.blockNumber, tc.blockHash)
	}

	for _, tc := range blockTable {
		h, err := o.GetBlockHash(tc.blockNumber)

		assert.Nil(err)
		assert.Equal(tc.blockHash, h)
	}
}

// Testing if GetBlockHash returns a copy
func Test_InMemory_GetBlockHash_ReturnsCopy(t *testing.T) {
	o := &oracle.InMemoryOracle{}
	number := uint64(42)

	o.AddBlock(number, common.BytesToHash([]byte{1}))

	h, _ := o.GetBlockHash(number)
	h.SetBytes([]byte{2})

	j, _ := o.GetBlockHash(number)

	assert.NotEqual(t, h, j)
	assert.NotEqual(t, common.Hash{}, j)
}

// Testing if GetBlockHash returns MissingBlockError
func Test_InMemory_GetBlockHash_MissingBlockError(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	number := uint64(132)

	_, err := o.GetBlockHash(number)

	var missingBlockError *oracle.MissingBlockError
	assert.ErrorAs(err, &missingBlockError)
	assert.Equal(number, missingBlockError.BlockNumber)

	assert.NotEmpty(missingBlockError.Msg)
	assert.NotEmpty(missingBlockError.Error())
}

// Testing AddTransaction on empty oracle
func Test_InMemory_AddTransaction(t *testing.T) {
	o := &oracle.InMemoryOracle{}

	for _, tc := range transactionTable {
		o.AddTransaction(tc.blockHash, tc.transactionIndex, tc.transactionHash)
	}
}

// Testing if AddTransaction makes copy
func Test_InMemory_AddTransaction_MakesCopy(t *testing.T) {
	o := &oracle.InMemoryOracle{}
	transactionIndex := uint64(42)
	blockHash := common.BytesToHash([]byte{1})
	transactionHash := common.BytesToHash([]byte{2})

	o.AddTransaction(blockHash, transactionIndex, transactionHash)
	transactionHash.SetBytes([]byte{3})

	h, _ := o.GetTransactionHash(blockHash, transactionIndex)

	assert.NotEqual(t, h, transactionHash)
}

// Testing TransactionHash
func Test_InMemory_GetTransactionHash(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	for _, tc := range transactionTable {
		o.AddTransaction(tc.blockHash, tc.transactionIndex, tc.transactionHash)
	}

	for _, tc := range transactionTable {
		h, err := o.GetTransactionHash(tc.blockHash, tc.transactionIndex)

		assert.Nil(err)
		assert.Equal(tc.transactionHash, h)
	}
}

// Testing if GetTransactionHash returns a copy
func Test_InMemory_GetTransactionHash_ReturnsCopy(t *testing.T) {
	o := &oracle.InMemoryOracle{}
	transactionIndex := uint64(42)
	blockHash := common.BytesToHash([]byte{1})
	transactionHash := common.BytesToHash([]byte{2})

	o.AddTransaction(blockHash, transactionIndex, transactionHash)
	h, _ := o.GetTransactionHash(blockHash, transactionIndex)
	h.SetBytes([]byte{3})

	j, _ := o.GetTransactionHash(blockHash, transactionIndex)

	assert.NotEqual(t, h, j)
	assert.NotEqual(t, common.Hash{}, j)
}

// Testing if GetTransactionHash returns MissingTransactionError
func Test_InMemory_GetTransactionHash_MissingTransactionError(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	hash := common.Hash{}
	index := uint64(132)

	_, err := o.GetTransactionHash(hash, index)

	var missingTransactionError *oracle.MissingTransactionError
	assert.ErrorAs(err, &missingTransactionError)
	assert.Equal(hash, missingTransactionError.BlockHash)
	assert.Equal(index, missingTransactionError.TransactionIndex)

	assert.NotEmpty(missingTransactionError.Msg)
	assert.NotEmpty(missingTransactionError.Error())
}

// Testing the API on empty transcript oracle
func Test_Transcript_Empty(t *testing.T) {
	assert := assert.New(t)
	to := &oracle.TranscriptOracle{}

	_, err := to.GetBlockHash(uint64(0))
	assert.NotNil(err)

	_, err = to.GetTransactionHash(common.Hash{}, uint64(0))
	assert.NotNil(err)

	transcript := to.GetTranscript()
	assert.Equal(0, len(transcript))
}

// Testing if SetOracle resets transcript
func Test_Transcript_SetOracle_Resets_Transcript(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	for _, tc := range blockTable {
		o.AddBlock(tc.blockNumber, tc.blockHash)
	}

	to := oracle.NewTranscriptOracle(o)

	for _, tc := range blockTranscriptTable {
		_, _ = to.GetBlockHash(tc.blockNumber)
	}

	assert.Equal(len(blockTranscriptTable), len(to.GetTranscript()))

	to.SetOracle(nil)

	assert.Equal(0, len(to.GetTranscript()))
}

// Testing if GetBlockHash returns MissingBlockError
func Test_Transcript_GetBlockHash_MissingBlockError(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}
	to := oracle.NewTranscriptOracle(o)
	number := uint64(1)

	_, err := to.GetBlockHash(number)

	var missingBlockError *oracle.MissingBlockError
	assert.ErrorAs(err, &missingBlockError)
	assert.Equal(number, missingBlockError.BlockNumber)

	assert.NotEmpty(missingBlockError.Msg)
	assert.NotEmpty(missingBlockError.Error())
}

// Testing if GetTransactionHash returns Missing TransactionError
func Test_Transcript_GetTransactionHash_MissingTransactionError(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}
	to := oracle.NewTranscriptOracle(o)
	hash := common.BytesToHash([]byte{1})
	index := uint64(1)

	_, err := to.GetTransactionHash(hash, index)

	var missingTransactionError *oracle.MissingTransactionError
	assert.ErrorAs(err, &missingTransactionError)
	assert.Equal(hash, missingTransactionError.BlockHash)
	assert.Equal(index, missingTransactionError.TransactionIndex)

	assert.NotEmpty(missingTransactionError.Msg)
	assert.NotEmpty(missingTransactionError.Error())
}

// Testing calling block transcript
func Test_Transcript_BlockTranscript(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	for _, tc := range blockTable {
		o.AddBlock(tc.blockNumber, tc.blockHash)
	}

	to := oracle.NewTranscriptOracle(o)

	for _, tc := range blockTranscriptTable {
		h, err := to.GetBlockHash(tc.blockNumber)
		assert.Equal(tc.blockHash, h)
		assert.Nil(err)
	}

	transcript := to.GetTranscript()
	assert.Equal(len(blockTranscriptTable), len(transcript))

	for i := range transcript {
		fact, ok := transcript[i].(oracle.BlockFact)
		assert.Truef(ok, "fact %v couldnt be cast to BlockFact", fact)
		tc := blockTranscriptTable[i]

		assert.Equal(tc.blockNumber, fact.BlockNumber)
		assert.Equal(tc.blockHash, fact.BlockHash)
	}
}

// Testing calling transaction transcript
func Test_Transcript_TransactionTranscript(t *testing.T) {
	assert := assert.New(t)
	o := &oracle.InMemoryOracle{}

	for _, tc := range transactionTable {
		o.AddTransaction(tc.blockHash, tc.transactionIndex, tc.transactionHash)
	}

	to := oracle.NewTranscriptOracle(o)

	for _, tc := range transactionTranscriptTable {
		h, err := to.GetTransactionHash(tc.blockHash, tc.transactionIndex)
		assert.Equal(tc.transactionHash, h)
		assert.Nil(err)
	}

	transcript := to.GetTranscript()
	assert.Equal(len(transactionTranscriptTable), len(transcript))

	for i := range transcript {
		fact, ok := transcript[i].(oracle.TransactionFact)
		assert.Truef(ok, "fact %v couldnt be cast to TransactionFact", fact)
		tc := transactionTranscriptTable[i]

		assert.Equal(tc.blockHash, fact.BlockHash)
		assert.Equal(tc.transactionIndex, fact.TransactionIndex)
		assert.Equal(tc.transactionHash, fact.TransactionHash)
	}
}

// Small benchmark for InMemory oracle
func Benchmark_InMemory(t *testing.B) {
	o := &oracle.InMemoryOracle{}

	for n := 0; n < t.N; n++ {
		for _, tc := range blockTable {
			o.AddBlock(tc.blockNumber, tc.blockHash)
		}

		for _, tc := range transactionTable {
			o.AddTransaction(tc.blockHash, tc.transactionIndex, tc.transactionHash)
		}
	}

	_result, _ = o.GetBlockHash(uint64(0))
}

// Small benchmark for Transcript oracle
func Benchmark_Transcript(t *testing.B) {
	o := &oracle.InMemoryOracle{}

	for _, tc := range blockTable {
		o.AddBlock(tc.blockNumber, tc.blockHash)
	}

	for _, tc := range transactionTable {
		o.AddTransaction(tc.blockHash, tc.transactionIndex, tc.transactionHash)
	}

	to := oracle.NewTranscriptOracle(o)

	r := common.Hash{}

	for n := 0; n < t.N; n++ {
		for _, tc := range blockTranscriptTable {
			r, _ = to.GetBlockHash(tc.blockNumber)
		}

		for _, tc := range transactionTranscriptTable {
			r, _ = to.GetTransactionHash(tc.blockHash, tc.transactionIndex)
		}

		transcript := to.GetTranscript()

		if len(transcript) > 0 {
			r = transcript[0].(oracle.BlockFact).BlockHash
		}
	}

	_result = r
}
