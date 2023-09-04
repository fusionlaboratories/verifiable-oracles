package flat

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
)

var encodingTable = []struct {
	b    []byte
	want Element
}{
	{[]byte{0x1}, field.NewElement(0x1)},
	{[]byte{0x11}, field.NewElement(0x11)},
	{[]byte{0x11, 0x22}, field.NewElement(0x11_22)},
	{[]byte{0x11, 0x22, 0x33}, field.NewElement(0x11_22_33)},
	{[]byte{0x11, 0x22, 0x33, 0x44}, field.NewElement(0x11_22_33_44)},
}

func TestEncoding_encodeBytes(t *testing.T) {
	for _, tc := range encodingTable {
		t.Run(fmt.Sprintf("%#v", tc.b), func(t *testing.T) {
			var e Element
			encodeBytes(tc.b, &e)

			assert.Equal(t, tc.want, e)

			// decoding
			b := make([]byte, len(tc.b))
			decodeBytes(&e, b)
			assert.Equal(t, tc.b, b, "Actual bytes %#v", e.Bytes())
		})
	}
}

func TestEncoding_Property(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MaxSize = _bytesInElement + 1
	properties := gopter.NewProperties(parameters)

	properties.Property("encode bytes", prop.ForAll(func(b []byte) (msg string) {
		var e Element
		var c = make([]byte, len(b))
		encodeBytes(b, &e)
		decodeBytes(&e, c)

		if !reflect.DeepEqual(b, c) {
			msg = fmt.Sprintf("%#v != %#v (%#v)", b, c, e.Bytes())
		}
		return
	}, gen.SliceOf(gen.UInt8())))

	properties.TestingRun(t)
}
