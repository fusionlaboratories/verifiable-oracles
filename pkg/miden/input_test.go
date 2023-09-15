package miden_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"github.com/qredo/verifiable-oracles/pkg/miden"
)

var _inputFileTable = map[string]struct {
	input miden.Input
	want  string
}{
	"zero inputFile": {
		input: miden.Input{},
		want:  `{"operand_stack":[]}`,
	},
	"non-empty operand stack": {
		input: miden.Input{
			OperandStack: field.Vector{field.One()},
		},
		want: `{"operand_stack":["1"]}`,
	},
	"non-empty advice stack": {
		input: miden.Input{
			AdviceStack: field.Vector{field.One()},
		},
		want: `{"advice_stack":["1"],"operand_stack":[]}`,
	},
	"both stacks non-empty": {
		input: miden.Input{
			OperandStack: field.Vector{field.One()},
			AdviceStack:  field.Vector{field.One()},
		},
		want: `{"advice_stack":["1"],"operand_stack":["1"]}`,
	},
}

func TestInputFileJsonMarshal(t *testing.T) {
	for name, tc := range _inputFileTable {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			j, err := json.Marshal(tc.input)

			assert.Nil(err)
			assert.Equal(tc.want, string(j))
		})
	}
}

func TestInputTestData(t *testing.T) {
	assert := assert.New(t)

	data, err := os.ReadFile("testdata/input.json")
	assert.Nil(err)

	var in miden.Input
	err = json.Unmarshal(data, &in)
	assert.Nil(err)
}
