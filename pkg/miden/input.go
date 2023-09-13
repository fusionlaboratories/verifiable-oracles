package miden

import (
	"encoding/json"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

// Input File for Miden
type Input struct {
	OperandStack field.Vector `json:"operand_stack"`
	AdviceStack  field.Vector `json:"advice_stack,omitempty"`
}

func marshalVector(v field.Vector) []string {
	r := make([]string, len(v))
	for i := range v {
		r[i] = v[i].String()
	}
	return r
}

// Need to explicitly implement json.Marshaler, as Miden expect expty stacks to
// be encoded as []
func (f Input) MarshalJSON() ([]byte, error) {
	data := map[string]any{}

	data["operand_stack"] = marshalVector(f.OperandStack)
	if len(f.AdviceStack) != 0 {
		data["advice_stack"] = marshalVector(f.AdviceStack)
	}

	return json.Marshal(data)
}

// Type Assertions
var _ json.Marshaler = (*Input)(nil)
var _ json.Marshaler = Input{}
