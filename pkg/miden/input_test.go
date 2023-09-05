package miden_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/qredo/verifiable-oracles/pkg/miden"
	"github.com/stretchr/testify/assert"
)

func TestInputFileJsonEncode(t *testing.T) {
	assert := assert.New(t)
	var i miden.InputFile

	j, err := json.Marshal(i)

	assert.Nil(err)
	assert.Equal(`{"operand_stack":null}`, string(j))
}

func TestTestData(t *testing.T) {
	assert := assert.New(t)

	_, err := os.ReadFile("testdata/input.json")
	assert.Nil(err)
}
