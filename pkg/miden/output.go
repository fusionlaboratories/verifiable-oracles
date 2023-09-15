package miden

import (
	"encoding/json"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

// Output holds the output of running a Miden VM program
type Output struct {
	Stack         field.Vector `json:"stack"`
	OverflowAddrs []string     `json:"overflow_addrs"`
}

func (f Output) MarshalJSON() ([]byte, error) {
	data := make(map[string]any, 2)

	data["stack"] = marshalVector(f.Stack)
	if len(f.OverflowAddrs) > 0 {
		data["overflow_addrs"] = f.OverflowAddrs
	} else {
		data["overflow_addrs"] = []string{}
	}

	return json.Marshal(data)
}

// Type Assertions
var _ json.Marshaler = (*Output)(nil)
var _ json.Marshaler = Output{}
