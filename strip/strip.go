package strip

import (
	"math"
	"math/rand"
	"sync"
)

type LEDStrip struct {
	Size int
	Buffer

	Width, Height int
	sync.Mutex
}

func New(size int) *LEDStrip {
	width := int(math.Sqrt(float64(size)))
	return &LEDStrip{
		Size:   size,
		Buffer: NewBuffer(size),
		Width:  width,
		Height: width,
	}
}

func (s *LEDStrip) SetPixel(i int, color []uint8) {
	if 0 <= i && i < len(s.Buffer) {
		copy(s.Buffer[i], color)
	}
}

func (s *LEDStrip) SetXY(x, y int, color LED) {
	s.SetPixel(y*s.Width+x, color)
}

func (s *LEDStrip) GetXY(x, y int) LED {
	return s.Buffer[y*s.Width+x]
}

func (s *LEDStrip) Randomize() {
	for i := range s.Buffer {
		s.Buffer[i] = []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	}
}
