package io

import field "github.com/qredo/verifiable-oracles/pkg/goldilocks"

type VectorReader interface {
	ReadVector(field.Vector) (int, error)
}

type VectorWriter interface {
	WriteVector(v field.Vector) (n int, err error)
}
