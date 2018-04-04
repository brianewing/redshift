package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type RainbowEffect struct {
	Size    int
	Speed   float64
	Reverse bool
	Blend   Blend

	wheel [][]uint8
}

func NewRainbowEffect() *RainbowEffect {
	return &RainbowEffect{
		Size:  100,
		Speed: 0.5,
		Blend: *NewBlend(),
	}
}

var granularity int = 60 // generate n times more colours for better transitions

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	steps := e.Size * granularity
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
