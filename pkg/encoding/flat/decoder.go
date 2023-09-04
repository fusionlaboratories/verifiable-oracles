package flat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/qredo/verifiable-oracles/pkg/io"
)

type Decoder struct {
	r io.VectorReader
}

func NewDecoder(r io.VectorReader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(x any) (int, error) {
	if x == nil {
		return 0, errors.New("decode nil")
	}

	return d.DecodeValue(reflect.ValueOf(x))
}

func (d *Decoder) DecodeValue(x reflect.Value) (int, error) {
	if !x.IsValid() {
		return 0, errors.New("decode invalid reflect.Value")
	}

	switch x.Kind() {
	default:
		return 0, fmt.Errorf("decode to unsupported kind %s", x.Kind().String())
	case reflect.Pointer:
		switch x.Elem().Kind() {
		default:
			return 0, fmt.Errorf("decode to unsupported type %s", x.Type().String())
		case reflect.Uint64:
			v := x.Interface().(*uint64)
			return 1, d.decodeUint64(v)
		case reflect.Uint32:
			v := x.Interface().(*uint32)
			return 1, d.decodeUint32(v)
		case reflect.Uint16:
			v := x.Interface().(*uint16)
			return 1, d.decodeUint16(v)
		case reflect.Uint8:
			v := x.Interface().(*uint8)
			return 1, d.decodeUint8(v)
		}
	}
}

func (d *Decoder) DecodeBytes(b []byte) (int, error) {
	l := len(b)
	m := l / _bytesInElement

	if (l % _bytesInElement) > 0 {
		m += 1
	}
	v := make(Vector, m)

	n, err := d.r.ReadVector(v)
	if err != nil {
		return 0, err
	}

	for i := 0; i < n; i++ {
		j := i * _bytesInElement
		src := &v[i]
		dst := b[j:min(l, j+_bytesInElement)]

		decodeBytes(src, dst)
	}

	return min(n*_bytesInElement, l), nil
}

func (d *Decoder) decodeNBytes(n int) ([]byte, error) {
	var b []byte = make([]byte, n)
	if m, err := d.DecodeBytes(b); err != nil {
		return nil, err
	} else if m < n {
		return nil, errors.New("not enough bytes")
	}
	return b, nil
}

func (d *Decoder) decodeUint8(r *uint8) (err error) {
	b, err := d.decodeNBytes(1)
	if err == nil {
		*r = b[0]
	}
	return
}

func (d *Decoder) decodeUint16(r *uint16) (err error) {
	b, err := d.decodeNBytes(2)
	if err == nil {
		*r = binary.BigEndian.Uint16(b)
	}
	return
}

func (d *Decoder) decodeUint32(r *uint32) (err error) {
	b, err := d.decodeNBytes(4)
	if err == nil {
		*r = binary.BigEndian.Uint32(b)
	}
	return
}

func (d *Decoder) decodeUint64(r *uint64) (err error) {
	b, err := d.decodeNBytes(8)
	if err == nil {
		*r = binary.BigEndian.Uint64(b)
	}
	return
}
