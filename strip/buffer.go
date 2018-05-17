package strip

import (
	"bytes"
	"strconv"
)

type Buffer [][]uint8

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

func (b *Buffer) MarshalJSON() ([]byte, error) {
	var tmp bytes.Buffer
	tmp.WriteRune('[')
	for i, led := range *b {
		if i != 0 {
			tmp.WriteRune(',')
		}
		tmp.WriteRune('[')
		for j, val := range led {
			if j != 0 {
				tmp.WriteRune(',')
			}
			tmp.WriteString(strconv.Itoa(int(val)))
		}
		tmp.WriteRune(']')
	}
	tmp.WriteRune(']')
	return tmp.Bytes(), nil
}
