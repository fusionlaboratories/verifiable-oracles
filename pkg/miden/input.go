package miden

import (
	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

// Input File for
type InputFile struct {
	OperandStack field.Vector `json:"operand_stack"`
	AdviceStack  field.Vector `json:"advice_stack,omitempty"`
}
