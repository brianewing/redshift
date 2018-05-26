package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type RainbowEffect struct {
	Size    int
	Depth   int
	Speed   float64
	Reverse bool
	Blend   Blend

	wheel [][]uint8
}

func NewRainbowEffect() *RainbowEffect {
	return &RainbowEffect{
		Size:  100,
		Depth: 60,
		Speed: 0.1,
		Blend: *NewBlend(),
	}
}

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	if e.Depth == 0 {
		e.Depth = 1
	}
	steps := e.Size * e.Depth

	if e.wheel == nil || len(e.wheel) != int(steps) {
		e.wheel = generateWheel(steps)
	}

	phase := round(CycleBetween(0, float64(len(e.wheel)), e.Speed))

	rotatedWheel := rotateBuffer(e.wheel, phase, e.Reverse)
	sampledWheel := sampleBuffer(rotatedWheel, int(e.Size))

	e.Blend.Buffer = sampledWheel
	e.Blend.Render(s)
}

func generateWheel(size int) [][]uint8 {
	wheel := make([][]uint8, size)

	for i := range wheel {
		hue := float64(i) / float64(size) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = []uint8{r, g, b}
	}

	return wheel
}
