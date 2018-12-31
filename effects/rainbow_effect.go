package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type RainbowEffect struct {
	Size    int
	Depth   int
	Speed   float64
	Blend   *Blend
	Reverse bool

	wheel strip.Buffer
}

func NewRainbowEffect() *RainbowEffect {
	return &RainbowEffect{
		Size:  100,
		Depth: 60,
		Speed: 0.1,
		Blend: NewBlend(),
	}
}

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	if e.Depth == 0 {
		e.Depth = 1
	}

	steps := e.Size * e.Depth

	if e.wheel == nil || len(e.wheel) != int(steps) {
		e.wheel = strip.Buffer(generateWheel(steps))
	}

	phase := round(CycleBetween(0, float64(len(e.wheel)), e.Speed))

	e.Blend.Reverse = e.Reverse
	e.Blend.Buffer = e.wheel.Rotate(phase, false).Sample(e.Size)
	e.Blend.Render(s)
}

func generateWheel(size int) strip.Buffer {
	wheel := make(strip.Buffer, size)

	for i := range wheel {
		hue := float64(i) / float64(size) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = strip.LED{r, g, b}
	}

	return wheel
}
