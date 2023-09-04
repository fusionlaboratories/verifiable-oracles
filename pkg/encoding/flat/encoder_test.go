package flat_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/elements"
	"github.com/qredo/verifiable-oracles/pkg/encoding/flat"
	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

// TODO: Add Benchmarks
type Element = field.Element
type Vector = field.Vector

var tableBytes = []struct {
	input []byte
	want  Vector
}{
	{[]byte{0x00}, Vector{field.NewElement(0x00)}},
	{[]byte{0x01}, Vector{field.NewElement(0x01)}},
	{[]byte{0x02}, Vector{field.NewElement(0x02)}},
	{[]byte{0x12, 0x34, 0x56, 0x78, 0x9A}, Vector{field.NewElement(0x12_34_56_78), field.NewElement(0x9A)}},
}

func TestEncodeNil(t *testing.T) {
	var e flat.Encoder
	n, err := e.Encode(nil)
	assert.Equal(t, 0, n)
	assert.Error(t, err)
}

func TestEncodeBytes(t *testing.T) {
	for _, tc := range tableBytes {

		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			var r elements.ElementBuffer
			e := flat.NewEncoder(&r)

			l, err := e.EncodeBytes(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, len(tc.input), l)
			assert.Equal(t, tc.want, r.Flush())
		})
	}
}

// Test if numbers are big-endian
func TestEndianess(t *testing.T) {
	var a elements.ElementBuffer
	var b elements.ElementBuffer
	var e = flat.NewEncoder(&a)
	var f = flat.NewEncoder(&b)

	e.Encode(0x12_34_56_78_9A)
	f.Encode(0x12)
	f.Encode(0x34_56_78_9A)

	assert.Equal(t, a.Flush(), b.Flush())
}
