package strip

import (
	"math/rand"
	"sync"
)

type LEDStrip struct {
	Size int
	Buffer [][]int

	Sync sync.Mutex
}

func New(size int) *LEDStrip {
	buffer := make([][]int, size)
	return &LEDStrip{Size: size, Buffer: buffer}
}

func (s *LEDStrip) Clear() {
	for i := range s.Buffer {
		s.Buffer[i] = []int{0, 0, 0}
	}
}

func (s *LEDStrip) SetPixel(i int, color []int) {
	if 0 <= i && i < len(s.Buffer) {
		copy(s.Buffer[i], color)
	}
}

func (s *LEDStrip) Randomize() {
	for i:= range s.Buffer {
		s.Buffer[i] = []int{rand.Intn(255), rand.Intn(255), rand.Intn(255)}
	}
}
