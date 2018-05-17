package strip

import (
	"math/rand"
	"sync"
)

type LEDStrip struct {
	Size int
	Buffer
	sync.Mutex
}

func New(size int) *LEDStrip {
	return &LEDStrip{Size: size, Buffer: NewBuffer(size)}
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
