package strip

import (
	"math/rand"
	"sync"
)

type LEDStrip struct {
	Size int
	Buffer [][]uint8

	sync.Mutex
}

func New(size int) *LEDStrip {
	return &LEDStrip{Size: size, Buffer: NewBuffer(size)}
}

func NewBuffer(size int) [][]uint8 {
	buffer := make([][]uint8, size)
	clearBuffer(buffer)
	return buffer
}

func clearBuffer(buffer [][]uint8) {
	for i := range buffer {
		buffer[i] = []uint8{0, 0, 0}
	}
}

func (s *LEDStrip) Clear() {
	clearBuffer(s.Buffer)
}

func (s *LEDStrip) SetPixel(i int, color []uint8) {
	if 0 <= i && i < len(s.Buffer) {
		copy(s.Buffer[i], color)
	}
}

func (s *LEDStrip) Randomize() {
	for i:= range s.Buffer {
		s.Buffer[i] = []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	}
}
