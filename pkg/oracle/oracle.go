package oracle

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// Custom error type for missing blocks
type MissingBlockError struct {
	BlockNumber *big.Int
	Msg         string
}

func (e *MissingBlockError) Error() string {
	return e.Msg
}

// Oracle Queries
type Oracle interface {
	GetBlockHash(number *big.Int) (common.Hash, error)
}

type InMemoryOracle struct {
	hashes map[uint64]common.Hash
	mu     sync.RWMutex
}

func NewInMemoryOracle() *InMemoryOracle {
	return &InMemoryOracle{hashes: make(map[uint64]common.Hash)}
}

func (o *InMemoryOracle) GetBlockHash(number *big.Int) (common.Hash, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// TODO: Add boundary check for Uint64
	n := number.Uint64()
	h, ok := o.hashes[n]
	if !ok {
		return common.Hash{}, &MissingBlockError{
			BlockNumber: number,
			Msg:         fmt.Sprintf("missing block %d", number),
		}
	}
	return h, nil
}

func (o *InMemoryOracle) AddBlock(blockNumber *big.Int, blockHash common.Hash) {
	o.mu.Lock()
	defer o.mu.Unlock()

	n := blockNumber.Uint64()

	o.hashes[n] = blockHash
}

type TranscriptOracle struct {
	inner Oracle
	mu    sync.RWMutex

	transcript []BlockFact
}

func NewTranscriptOracle(inner Oracle) *TranscriptOracle {
	return &TranscriptOracle{inner: inner}
}

func (to *TranscriptOracle) GetBlockHash(number *big.Int) (common.Hash, error) {
	hash, err := to.inner.GetBlockHash(number)
	if err != nil {
		// TODO: Consider Wrapping the Error
		return hash, err
	}

	to.mu.Lock()
	defer to.mu.Unlock()

	to.transcript = append(to.transcript, BlockFact{Number: number, Hash: hash})

	return hash, nil
}

func (to *TranscriptOracle) GetTranscript() []BlockFact {
	to.mu.RLock()
	defer to.mu.RUnlock()

	transcript := make([]BlockFact, len(to.transcript))
	copy(transcript, to.transcript)
	return transcript
}

var _ Oracle = (*TranscriptOracle)(nil)
var _ Oracle = (*InMemoryOracle)(nil)

// TODO: Encode as JSON
type BlockFact struct {
	// TODO: Is this immutable?
	Number *big.Int
	Hash   common.Hash
}

/*
// Transcript Entry for getting blocks

// TODO: Transcript Entry for getting transactions
type TransactionFact struct {
}

// TODO: What kind of interface transcript entry needs to have
type Fact interface {
	// TODO: print as string
	// TODO: encode as JSON
	// TODO: encode as Field elements
}
*/
