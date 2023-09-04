package elements

import (
	"io"

	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	ffio "github.com/qredo/verifiable-oracles/pkg/io"
)

// TODO: Implement more API of Buffer
type ElementBuffer struct {
	s      field.Vector
	offset int
}

func NewElementBuffer(s field.Vector) *ElementBuffer {
	return &ElementBuffer{s: s}
}

func (eb *ElementBuffer) Count() int {
	return len(eb.s) - eb.offset
}

func (eb *ElementBuffer) empty() bool {
	return eb.Count() <= 0
}

func (eb *ElementBuffer) ReadElement() (field.Element, error) {
	if eb.empty() {
		eb.Reset()
		return field.Element{}, io.EOF
	}

	e := eb.s[eb.offset]
	eb.offset++

	return e, nil
}

func (eb *ElementBuffer) ReadVector(v field.Vector) (int, error) {
	if eb.empty() {
		eb.Reset()
		if len(v) == 0 {
			return 0, nil
		}

		return 0, io.EOF
	}

	n := copy(v, eb.s[eb.offset:])
	eb.offset += n

	return n, nil
}

func (eb *ElementBuffer) WriteElement(e field.Element) error {
	eb.s = append(eb.s, e)
	return nil
}

func (eb *ElementBuffer) WriteVector(v field.Vector) (int, error) {
	eb.s = append(eb.s, v...)
	return len(v), nil
}

func (eb *ElementBuffer) Vector() field.Vector {
	return eb.s[eb.offset:]
}

func (eb *ElementBuffer) Flush() field.Vector {
	v := eb.Vector()
	eb.Reset()
	return v
}

func (eb *ElementBuffer) Reset() {
	eb.s = nil
	eb.offset = 0
}

var _ ffio.ElementReader = (*ElementBuffer)(nil)
var _ ffio.VectorReader = (*ElementBuffer)(nil)
var _ ffio.ElementWriter = (*ElementBuffer)(nil)
var _ ffio.VectorWriter = (*ElementBuffer)(nil)
