package miden

import (
	"encoding/json"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

// Input File for
type InputFile struct {
	OperandStack field.Vector `json:"operand_stack"`
	AdviceStack  field.Vector `json:"advice_stack,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (f InputFile) MarshalJSON() ([]byte, error) {
	data := map[string]any{}

	operand_stack := make([]string, len(f.OperandStack))
	for i := range operand_stack {
		operand_stack[i] = f.OperandStack[i].String()
	}
	data["operand_stack"] = operand_stack

	if len(f.AdviceStack) > 0 {
		data["advice_stack"] = f.AdviceStack
	}

	return json.Marshal(data)
}

var _ json.Marshaler = (*InputFile)(nil)
var _ json.Marshaler = InputFile{}
