package io

import field "github.com/qredo/verifiable-oracles/pkg/goldilocks"

type ElementReader interface {
	ReadElement() (field.Element, error)
}

type ElementWriter interface {
	WriteElement(e field.Element) error
}
