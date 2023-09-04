package flat_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/elements"
	"github.com/qredo/verifiable-oracles/pkg/encoding/flat"
	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

var reflectTable = []struct {
	input any
	want  Vector
}{
	{uint64(0x0), Vector{{}, {}}},
	{uint64(0x1), Vector{{}, field.One()}},
	{uint64(0x3), Vector{{}, field.NewElement(0x3)}},
	{uint64(0x5), Vector{{}, field.NewElement(0x5)}},
	{uint64(0x5), Vector{{}, field.NewElement(0x5)}},
	{uint64(0x7), Vector{{}, field.NewElement(0x7)}},
	{uint64(0x12_34), Vector{{}, field.NewElement(0x12_34)}},
	{uint64(0x24_68), Vector{{}, field.NewElement(0x24_68)}},
	{uint64(0x12_34_56_78_9A), Vector{field.NewElement(0x12), field.NewElement(0x34_56_78_9A)}},
	{uint32(0x0), Vector{{}}},
	{uint32(0x1), Vector{field.One()}},
	{uint32(0x2), Vector{field.NewElement(0x2)}},
	{uint32(0x3), Vector{field.NewElement(0x3)}},
	{uint16(0x0), Vector{{}}},
	{uint16(0x1), Vector{field.One()}},
	{uint16(0x2), Vector{field.NewElement(0x2)}},
	{uint16(0x3), Vector{field.NewElement(0x3)}},
	{uint8(0x0), Vector{{}}},
	{uint8(0x1), Vector{field.One()}},
	{uint8(0x2), Vector{field.NewElement(0x2)}},
	{uint8(0x3), Vector{field.NewElement(0x3)}},
}

// TODO: Split these in a separate tests
func TestEncodeDecode(t *testing.T) {
	for _, tc := range reflectTable {
		name := fmt.Sprintf("%#v", tc.input)
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)

				input  any = tc.input
				output     = reflect.New(reflect.TypeOf(input))

				buf elements.ElementBuffer
				dec = flat.NewDecoder(&buf)
				enc = flat.NewEncoder(&buf)

				n   int
				err error
			)

			n, err = enc.Encode(input)
			assert.Equal(1, n)
			assert.Nil(err)

			assert.Equal(tc.want, buf.Vector())

			n, err = dec.Decode(output.Interface())
			assert.Equal(1, n)
			assert.Nil(err)

			assert.Equal(input, output.Elem().Interface())
		})
	}
}
