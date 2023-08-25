package prover

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

// A simple prover/verifier using encoding/gob encoder, the proof.
type GobProver struct {
}

func (*GobProver) Prove(transcript Transcript) (Proof, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(transcript); err != nil {
		return nil, fmt.Errorf("prover.Prove: %w", err)
	}

	return buf.Bytes(), nil
}

// TODO: Would just error return value suffice?
func (*GobProver) Verify(transcript Transcript, proof Proof) (bool, error) {
	reader := bytes.NewReader(proof)
	dec := gob.NewDecoder(reader)
	var decodedTranscript Transcript

	if err := dec.Decode(&decodedTranscript); err != nil {
		return false, fmt.Errorf("prover.Verify: %w", err)
	}

	return reflect.DeepEqual(transcript, decodedTranscript), nil
}

// We only provide some restricted set of provers and verifiers
var _ TranscriptProver = (*GobProver)(nil)
var _ TranscriptVerifier = (*GobProver)(nil)
