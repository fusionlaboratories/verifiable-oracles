package miden_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"github.com/qredo/verifiable-oracles/pkg/miden"
)

func TestInputFileJsonEncode(t *testing.T) {
	assert := assert.New(t)
	var i miden.InputFile

	j, err := json.Marshal(i)

	assert.Nil(err)
	assert.Equal(`{"operand_stack":[]}`, string(j))
}

func TestInputFileJsonEncode1(t *testing.T) {
	assert := assert.New(t)
	f := miden.InputFile{
		OperandStack: field.Vector{field.One()},
	}

	j, err := json.Marshal(f)

	assert.Nil(err)
	assert.Equal(`{"operand_stack":["1"]}`, string(j))
}

func TestTestData(t *testing.T) {
	assert := assert.New(t)

	_, err := os.ReadFile("testdata/input.json")
	assert.Nil(err)
}
