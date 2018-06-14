package strip

type Buffer []LED

func NewBuffer(size int) Buffer {
	buffer := make(Buffer, size)
	buffer.Clear()
	return buffer
}

func (b Buffer) Clear() {
	for i := range b {
		b[i] = []uint8{0, 0, 0}
	}
}

// returns a new slice containing the data in buffer rotated by n
func (b Buffer) Rotate(n int, reverse bool) Buffer {
	if reverse {
		head, tail := b[0:n], b[n:]
		return append(tail, head...)
	} else {
		head, tail := b[:len(b)-n], b[len(b)-n:]
		return append(tail, head...)
	}
}

// returns a subset of buffer (n evenly-spaced elements)
func (b Buffer) Sample(n int) Buffer {
	subset := make(Buffer, n)
	if n > 0 {
		step := len(b) / n
		for i := 0; i < n; i++ {
			subset[i] = b[i*step]
		}
	}
	return subset
}

func (b Buffer) MarshalBytes() []byte {
	bytes := make([]byte, len(b) * 3)
	for i, led := range b {
		y := i * 3
		bytes[y] = led[0]
		bytes[y+1] = led[1]
		bytes[y+2] = led[2]
	}
	return bytes
}

func (b Buffer) UnmarshalBytes(bytes []byte) {
	for i, val := range bytes {
		if len(b) == i / 3 {
			break
		}
		b[i / 3][i % 3] = val
	}
}
