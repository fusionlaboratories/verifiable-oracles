package prover

// As Proof is the result type of Prove, it has to be fixed
type Proof []byte

type Prover[T any] interface {
	Prove(transcript T) (Proof, error)
}

type Verifier[T any] interface {
	Verify(transcript T, proof Proof) (bool, error)
}
