package prover

import "github.com/ethereum/go-ethereum/common"

type Transcript struct {
	BlockNumber      uint64      `json:"blockNumber"`
	BlockHash        common.Hash `json:"blockHash"`
	TransactionIndex uint64      `json:"transactionIndex"`
	TransactionHash  common.Hash `json:"transactionHash"`
}

type TranscriptProver = Prover[Transcript]
type TranscriptVerifier = Verifier[Transcript]
