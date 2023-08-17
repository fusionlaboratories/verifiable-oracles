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

// Oracle Queries
type Oracle interface {
	GetBlockHash(number uint64) (common.Hash, error)
}

type InMemoryOracle struct {
	hashes map[uint64]common.Hash
	mu     sync.RWMutex
}

func NewInMemoryOracle() *InMemoryOracle {
	return &InMemoryOracle{hashes: make(map[uint64]common.Hash)}
}

func (o *InMemoryOracle) GetBlockHash(number uint64) (common.Hash, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	h, ok := o.hashes[number]
	if !ok {
		return common.Hash{}, &MissingBlockError{
			BlockNumber: number,
			Msg:         fmt.Sprintf("missing block %d", number),
		}
	}
	return h, nil
}

func (o *InMemoryOracle) AddBlock(blockNumber uint64, blockHash common.Hash) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.hashes[blockNumber] = blockHash
}

type TranscriptOracle struct {
	inner Oracle
	mu    sync.RWMutex

	transcript []BlockFact
}

func NewTranscriptOracle(inner Oracle) *TranscriptOracle {
	return &TranscriptOracle{inner: inner}
}

func (to *TranscriptOracle) GetBlockHash(number uint64) (common.Hash, error) {
	hash, err := to.inner.GetBlockHash(number)
	if err != nil {
		// TODO: Consider Wrapping the Error
		return hash, err
	}

	to.mu.Lock()
	defer to.mu.Unlock()

	to.transcript = append(to.transcript, BlockFact{BlockNumber: number, BlockHash: hash})

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

type BlockFact struct {
	BlockNumber uint64      `json:"blockNumber"`
	BlockHash   common.Hash `json:"blockHash"`
}
