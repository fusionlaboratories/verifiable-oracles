package flat_test

import (
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/qredo/verifiable-oracles/pkg/encoding/flat"

	"github.com/qredo/verifiable-oracles/pkg/elements"
)

// Property based test of encode/decode
func TestRoundTrip(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 10000

	properties := gopter.NewProperties(params)

	properties.Property("uint8", prop.ForAll(func(v uint8) bool {
		var buff elements.ElementBuffer
		enc := flat.NewEncoder(&buff)
		dec := flat.NewDecoder(&buff)
		var w uint8

		enc.Encode(v)
		dec.Decode(&w)

		return v == w
	}, gen.UInt8()))

	properties.Property("uint16", prop.ForAll(func(v uint16) bool {
		var buff elements.ElementBuffer
		enc := flat.NewEncoder(&buff)
		dec := flat.NewDecoder(&buff)
		var w uint16

		enc.Encode(v)
		dec.Decode(&w)

		return v == w

	}, gen.UInt16()))

	properties.Property("uint32", prop.ForAll(func(v uint32) bool {
		var buff elements.ElementBuffer
		enc := flat.NewEncoder(&buff)
		dec := flat.NewDecoder(&buff)
		var w uint32

		enc.Encode(v)
		dec.Decode(&w)

		return v == w

	}, gen.UInt32()))

	properties.Property("uint64", prop.ForAll(func(v uint64) bool {
		var buff elements.ElementBuffer
		enc := flat.NewEncoder(&buff)
		dec := flat.NewDecoder(&buff)
		var w uint64

		enc.Encode(v)
		dec.Decode(&w)

		return v == w

	}, gen.UInt64()))

	properties.Property("[]byte", prop.ForAll(func(v []byte) bool {
		var buff elements.ElementBuffer
		enc := flat.NewEncoder(&buff)
		dec := flat.NewDecoder(&buff)
		var w = make([]byte, len(v))

		enc.EncodeBytes(v)
		dec.DecodeBytes(w)

		return reflect.DeepEqual(v, w)

	}, gen.SliceOf(gen.UInt8())))

	properties.TestingRun(t)
}
