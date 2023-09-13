package miden_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"github.com/qredo/verifiable-oracles/pkg/miden"
)

var _outputTable = map[string]struct {
	data string
	want miden.Output
}{
	"empty": {
		data: "{}",
	},
	"non empty stack": {
		data: `{"stack": ["1", "2", "3"], "overflow_addr": []}`,
		want: miden.Output{
			Stack: field.Vector{field.NewElement(1), field.NewElement(2), field.NewElement(3)},
		},
	},
}

func TestOutputJsonUnmarshal(t *testing.T) {
	for name, tc := range _outputTable {
		t.Run(name, func(t *testing.T) {
			data := []byte(tc.data)
			var out miden.Output

			err := json.Unmarshal(data, &out)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, out)
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
