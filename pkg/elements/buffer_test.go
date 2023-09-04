package elements_test

import (
	"testing"

	"github.com/qredo/verifiable-oracles/pkg/elements"
	field "github.com/qredo/verifiable-oracles/pkg/goldilocks"
	"github.com/qredo/verifiable-oracles/pkg/io"
)

type Element = field.Element
type Vector = field.Vector

func readElement(elems []Element, r io.ElementReader) {
	for i := range elems {
		elems[i], _ = r.ReadElement()
	}
}

func readVector(elems []Element, r io.VectorReader) {
	r.ReadVector(elems)
}

func writeElement(elems []Element, w io.ElementWriter) {
	for _, e := range elems {
		w.WriteElement(e)
	}
}

func writeVector(elems []Element, w io.VectorWriter) {
	w.WriteVector(elems)
}

func genInput(count int) []Element {
	input := make([]Element, count)

	for i := 0; i < count; i++ {
		input[i] = field.NewElement(uint64(i))
	}

	return input
}

func benchElementWriter(inputSize int, b *testing.B) {
	input := genInput(inputSize)

	for i := 0; i < b.N; i++ {
		var b elements.ElementBuffer
		writeElement(input, &b)
	}
}

func BenchmarkElementBuffer_ElementWriter_10(b *testing.B) {
	benchElementWriter(10, b)
}

func BenchmarkElementBuffer_ElementWriter_100(b *testing.B) {
	benchElementWriter(100, b)
}

func BenchmarkElementBuffer_ElementWriter_1000(b *testing.B) {
	benchElementWriter(1000, b)
}

func benchVectorWriter(inputSize int, b *testing.B) {
	input := genInput(inputSize)

	for i := 0; i < b.N; i++ {
		var b elements.ElementBuffer
		writeVector(input, &b)
	}
}

func BenchmarkElementBuffer_VectorWriter_10(b *testing.B) {
	benchVectorWriter(10, b)
}

func BenchmarkElementBuffer_VectorWriter_100(b *testing.B) {
	benchVectorWriter(100, b)
}

func BenchmarkElementBuffer_VectorWriter_1000(b *testing.B) {
	benchVectorWriter(1000, b)
}

func benchElementReader(inputSize int, b *testing.B) {
	input := genInput(inputSize)
	output := make([]Element, inputSize)

	for i := 0; i < b.N; i++ {
		b := elements.NewElementBuffer(input)
		readElement(output, b)
	}
}

func BenchmarkElementBuffer_ElementReader_10(b *testing.B) {
	benchElementReader(10, b)
}

func BenchmarkElementBuffer_ElementReader_100(b *testing.B) {
	benchElementReader(100, b)
}

func BenchmarkElementBuffer_ElementReader_1000(b *testing.B) {
	benchElementReader(1000, b)
}

func benchVectorReader(inputSize int, b *testing.B) {
	input := genInput(inputSize)
	output := make([]Element, inputSize)

	for i := 0; i < b.N; i++ {
		b := elements.NewElementBuffer(input)
		readVector(output, b)
	}

}

func BenchmarkElementBuffer_VectorReader_10(b *testing.B) {
	benchVectorReader(10, b)
}

func BenchmarkElementBuffer_VectorReader_100(b *testing.B) {
	benchVectorReader(100, b)
}

func BenchmarkElementBuffer_VectorReader_1000(b *testing.B) {
	benchVectorReader(1000, b)
}
