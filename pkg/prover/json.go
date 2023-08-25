package prover

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// A simple prover where the proof is JSON encoding of the transcript

type JsonProver struct {
}

// Prove implements Prover.
func (*JsonProver) Prove(transcript Transcript) (Proof, error) {
	b, err := json.Marshal(transcript)

	if err != nil {
		return nil, fmt.Errorf("prover.Prove: %w", err)
	}

	return b, nil
}

// Verify implements Verifier.
func (*JsonProver) Verify(transcript Transcript, proof Proof) (bool, error) {
	var decodedTranscript Transcript

	if err := json.Unmarshal(proof, &decodedTranscript); err != nil {
		return false, fmt.Errorf("prover.Verify: %w", err)
	}

	return reflect.DeepEqual(transcript, decodedTranscript), nil
}

var _ TranscriptProver = (*JsonProver)(nil)
var _ TranscriptVerifier = (*JsonProver)(nil)
