package flat

const _bytesInElement = 4

// Encode bytes in element.  b should have at most _bytesInElement bytes,
// whereas e cannot be nil.
func encodeBytes(b []byte, e *Element) {
	e.SetBytes(b)
}

// Decode bytes encoded in Field element.  Notice that since internal
// representation is big endian, then bytes are aligned to the right.
func decodeBytes(e *Element, b []byte) {
	bytes := e.Bytes()

	// Bytes are aligned to the right.
	start := 8 - len(b)
	copy(b, bytes[start:])
}
