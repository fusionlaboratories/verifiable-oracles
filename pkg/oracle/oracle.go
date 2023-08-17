package oracle

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// Custom error type for missing blocks
type MissingBlockError struct {
	BlockNumber uint64
	Msg         string
}

func (e *MissingBlockError) Error() string {
	return e.Msg
}

type MissingTransactionError struct {
	BlockHash        common.Hash
	TransactionIndex uint64
	Msg              string
}

func (e *MissingTransactionError) Error() string {
	return e.Msg
}

var _ error = (*MissingBlockError)(nil)
var _ error = (*MissingTransactionError)(nil)

type Oracle interface {
	// Queries
	GetBlockHash(blockNumber uint64) (common.Hash, error)
	GetTransactionHash(blockHash common.Hash, transactionIndex uint64) (common.Hash, error)
}

type txKey struct {
	blockHash        common.Hash
	transactionIndex uint64
}

type InMemoryOracle struct {
	hashes       map[uint64]common.Hash
	numbers      map[common.Hash]uint64
	transactions map[txKey]common.Hash

	mu sync.RWMutex
}

func (o *InMemoryOracle) GetBlockHash(number uint64) (common.Hash, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	h, ok := o.hashes[number]
	if !ok {
		return common.Hash{}, &MissingBlockError{
			BlockNumber: number,
			Msg:         fmt.Sprintf("InMemoryOracle: missing block %d", number),
		}
	}
	return h, nil
}

func (o *InMemoryOracle) GetTransactionHash(blockHash common.Hash, transactionIndex uint64) (common.Hash, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	key := txKey{blockHash: blockHash, transactionIndex: transactionIndex}
	h, ok := o.transactions[key]
	if !ok {
		return common.Hash{}, &MissingTransactionError{
			BlockHash:        blockHash,
			TransactionIndex: transactionIndex,
			Msg:              fmt.Sprintf("InMemoryOracle: missing transaction %v %d", blockHash, transactionIndex),
		}
	}
	return h, nil
}

func (o *InMemoryOracle) AddBlock(blockNumber uint64, blockHash common.Hash) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.hashes == nil {
		o.hashes = map[uint64]common.Hash{}
	}
	o.hashes[blockNumber] = blockHash

	if o.numbers == nil {
		o.numbers = map[common.Hash]uint64{}
	}
	o.numbers[blockHash] = blockNumber
}

func (o *InMemoryOracle) AddTransaction(blockHash common.Hash, transactionIndex uint64, transactionHash common.Hash) {
	o.mu.Lock()
	defer o.mu.Unlock()

	key := txKey{blockHash: blockHash, transactionIndex: transactionIndex}

	if o.transactions == nil {
		o.transactions = map[txKey]common.Hash{}
	}
	o.transactions[key] = transactionHash
}

type TranscriptOracle struct {
	inner Oracle
	mu    sync.RWMutex

	transcript []any
}

func NewTranscriptOracle(inner Oracle) *TranscriptOracle {
	return &TranscriptOracle{inner: inner}
}

func (to *TranscriptOracle) SetOracle(inner Oracle) *TranscriptOracle {
	to.mu.Lock()
	defer to.mu.Unlock()

	to.inner = inner
	to.transcript = nil

	return to
}

func (to *TranscriptOracle) GetBlockHash(blockNumber uint64) (common.Hash, error) {
	to.mu.Lock()
	defer to.mu.Unlock()

	if to.inner == nil {
		return common.Hash{}, &MissingBlockError{
			BlockNumber: blockNumber,
			Msg:         fmt.Sprintf("TranscriptOracle: missing block %d", blockNumber),
		}
	}

	blockHash, err := to.inner.GetBlockHash(blockNumber)
	if err != nil {
		return blockHash, fmt.Errorf("TranscriptOracle.GetBlockHash: %w", err)
	}

	to.transcript = append(to.transcript, BlockFact{BlockNumber: blockNumber, BlockHash: blockHash})

	return blockHash, nil
}

func (to *TranscriptOracle) GetTransactionHash(blockHash common.Hash, transactionIndex uint64) (common.Hash, error) {
	to.mu.Lock()
	defer to.mu.Unlock()

	if to.inner == nil {
		return common.Hash{}, &MissingTransactionError{
			BlockHash:        blockHash,
			TransactionIndex: transactionIndex,
			Msg:              fmt.Sprintf("TranscriptOracle: missing transaction %v %d", blockHash, transactionIndex),
		}
	}

	transactionHash, err := to.inner.GetTransactionHash(blockHash, transactionIndex)
	if err != nil {
		return transactionHash, fmt.Errorf("TranscriptOracle.GetTransactionHash: %w", err)
	}

	to.transcript = append(to.transcript, TransactionFact{BlockHash: blockHash, TransactionIndex: transactionIndex, TransactionHash: transactionHash})

	return transactionHash, nil
}

func (to *TranscriptOracle) GetTranscript() Transcript {
	to.mu.RLock()
	defer to.mu.RUnlock()

	transcript := make(Transcript, len(to.transcript))
	copy(transcript, to.transcript)
	return transcript
}

var _ Oracle = (*TranscriptOracle)(nil)
var _ Oracle = (*InMemoryOracle)(nil)

type Transcript []any

type BlockFact struct {
	BlockNumber uint64      `json:"blockNumber"`
	BlockHash   common.Hash `json:"blockHash"`
}

type TransactionFact struct {
	BlockHash        common.Hash `json:"blockHash"`
	TransactionIndex uint64      `json:"transactionIndex"`
	TransactionHash  common.Hash `json:"transactionHash"`
}
