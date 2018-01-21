package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type RainbowEffect struct {
	Size    uint
	Speed   float64
	Reverse bool

	wheel   [][]uint8
}

var multiple uint = 60 // generate n times more colours for better transitions

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	size := e.getSize() * multiple
	if e.wheel == nil || len(e.wheel) != int(size) {
		e.wheel = generateWheel(size)
	}

	phase := round(CycleBetween(0, float64(len(e.wheel)), e.Speed))
	rotatedWheel := rotateBuffer(e.wheel, phase, e.Reverse)
	sampledWheel := sampleBuffer(rotatedWheel, int(e.getSize()))

	(&Buffer{Buffer: sampledWheel}).Render(s)
}

func (e *RainbowEffect) getSize() uint {
	if e.Size == 0 {
		e.Size = 150
	}
	return e.Size
}

func generateWheel(size uint) [][]uint8 {
	wheel := make([][]uint8, size)
	for i := range wheel {
		hue := float64(i) / float64(size) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = []uint8{r, g, b}
	}
	return wheel
}

