package flat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/qredo/verifiable-oracles/pkg/io"
)

// Encoder encodes data into field elements
type Encoder struct {
	w io.VectorWriter
}

func NewEncoder(w io.VectorWriter) *Encoder {
	e := &Encoder{w: w}
	return e
}

func (e *Encoder) Encode(x any) (int, error) {
	if x == nil {
		return 0, errors.New("encode nil")
	}

	return e.EncodeValue(reflect.ValueOf(x))
}

func (e *Encoder) EncodeValue(x reflect.Value) (int, error) {
	if !x.IsValid() {
		return 0, errors.New("invalid reflect.Value")

	}

	switch x.Kind() {
	default:
		return 0, fmt.Errorf("uknown Kind %s", x.Kind().String())
	case reflect.Uint64:
		v := x.Interface().(uint64)
		return 1, e.encodeUint64(v)
	case reflect.Uint32:
		v := uint32(x.Uint())
		return 1, e.encodeUint32(v)
	case reflect.Uint16:
		v := uint16(x.Uint())
		return 1, e.encodeUint16(v)
	case reflect.Uint8:
		v := uint8(x.Uint())
		return 1, e.encodeUint8(v)
	}
}

// Note that there will be padding added at the end
func (e *Encoder) EncodeBytes(b []byte) (int, error) {
	l := len(b)
	// TODO: Perform a single allocation
	var s Vector

	for i := 0; i < l; i += _bytesInElement {
		var e Element
		end := min(i+_bytesInElement, l)
		encodeBytes(b, &e)
		e.SetBytes(b[i:end])
		s = append(s, e)
	}

	_, err := e.append(s...)
	return l, err
}

func (e *Encoder) append(elems ...Element) (int, error) {
	return e.w.WriteVector(elems)
}

// TODO: Change the API of these functions to at least return an error
func (e *Encoder) encodeUint8(x uint8) error {
	var b []byte = []byte{x}
	_, err := e.EncodeBytes(b)
	return err
}

func (e *Encoder) encodeUint16(x uint16) error {
	var b []byte = make([]byte, 2)

	binary.BigEndian.PutUint16(b, x)
	_, err := e.EncodeBytes(b)
	return err
}

func (e *Encoder) encodeUint32(x uint32) error {
	var b []byte = make([]byte, 4)

	binary.BigEndian.PutUint32(b, x)
	_, err := e.EncodeBytes(b)
	return err
}

func (e *Encoder) encodeUint64(x uint64) error {
	var b []byte = make([]byte, 8)

	binary.BigEndian.PutUint64(b, x)
	_, err := e.EncodeBytes(b)
	return err
}

func (e *Encoder) EncodeElement(x Element) error {
	_, err := e.append(x)
	return err
}

func (e *Encoder) EncodeVector(x Vector) (int, error) {
	return e.append(x...)
}
