package effects

import (
	"github.com/brianewing/redshift/strip"
	opensimplex "github.com/ojrac/opensimplex-go"

	"math/rand"
)

type SimplexNoise struct {
	Width, Height int

	StepX, StepY     float64
	offsetX, offsetY float64

	noise opensimplex.Noise
}

func (e *SimplexNoise) Init() {
	e.noise = opensimplex.NewNormalized(rand.Int63())
}

func (e *SimplexNoise) Render(s *strip.LEDStrip) {
	e.offsetX += (e.StepX / 100)
	e.offsetY += (e.StepY / 100)

	for i, v := range e.Noise() {
		if i >= len(s.Buffer) {
			break
		}
		x := uint8(v * 255)
		copy(s.Buffer[i], strip.LED{x, x, x})
	}
}

func (e *SimplexNoise) Noise() (frame []float64) {
	for y := float64(0); y < float64(e.Width); y++ {
		for x := float64(0); x < float64(e.Height); x++ {
			val := e.noise.Eval2(e.offsetX+x, e.offsetY+y)
			frame = append(frame, val)
		}
	}
	return
}
