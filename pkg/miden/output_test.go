package miden_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"github.com/qredo/verifiable-oracles/pkg/miden"
)

var _outputUnmarshalTable = map[string]struct {
	data string
	want miden.Output
}{
	"empty": {
		data: `{}`,
	},
	"non empty stack": {
		data: `{"stack":["1", "2", "3"]}`,
		want: miden.Output{
			Stack: field.Vector{field.NewElement(1), field.NewElement(2), field.NewElement(3)},
		},
	},
	"overflow addresses": {
		data: `{"overflow_addrs":[], "stack":["1", "2", "3"]}`,
		want: miden.Output{
			Stack:         field.Vector{field.NewElement(1), field.NewElement(2), field.NewElement(3)},
			OverflowAddrs: []string{},
		},
	},
}

var _outputMarshalTable = map[string]struct {
	data miden.Output
	want string
}{
	"empty": {
		want: `{"overflow_addrs":[],"stack":[]}`,
	},
	"non empty stack": {
		data: miden.Output{
			Stack: field.Vector{field.NewElement(1), field.NewElement(2), field.NewElement(3)},
		},
		want: `{"overflow_addrs":[],"stack":["1","2","3"]}`,
	},
}

func TestOutputJsonUnmarshal(t *testing.T) {
	for name, tc := range _outputUnmarshalTable {
		t.Run(name, func(t *testing.T) {
			data := []byte(tc.data)
			var out miden.Output

			err := json.Unmarshal(data, &out)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestOutputJsonMarshal(t *testing.T) {
	for name, tc := range _outputMarshalTable {
		t.Run(name, func(t *testing.T) {

			data, err := json.Marshal(tc.data)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, string(data))
		})
	}
}

func TestOutputTestData(t *testing.T) {
	assert := assert.New(t)
	var out miden.Output

	data, err := os.ReadFile("testdata/output.json")
	assert.Nil(err)

	err = json.Unmarshal(data, &out)
	assert.Nil(err)
}
