package effects

import (
	"math/rand"

	"github.com/brianewing/redshift/strip"
)

type LarsonEffect struct {
	Color    strip.LED
	Speed    float64
	Position int
	Width    int
	Bounce   bool
}

func NewLarsonEffect() *LarsonEffect {
	return &LarsonEffect{
		Width: 2,
		Color: strip.LED{0, 0, 0},
		Speed: 0.1,
	}
}

func (e *LarsonEffect) Render(s *strip.LEDStrip) {
	if e.Speed != 0 {
		var fn TimingFunction
		speed := e.Speed
		if e.Bounce {
			fn = OscillateBetween
			speed /= 2
		} else {
			fn = CycleBetween
		}
		e.Position = round(fn(0, float64(s.Size-e.Width), speed))
	}

	color := e.getColor()
	for i := 0; i < e.Width; i++ {
		s.SetPixel(e.Position+i, color)
	}
}

func (e *LarsonEffect) getColor() []uint8 {
	if len(e.Color) == 0 {
		return []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	}

	return e.Color
}
