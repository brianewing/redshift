package effects

import (
	"github.com/brianewing/redshift/strip"
	"math/rand"
)

type LarsonEffect struct {
	Color    []uint8
	Position int
	Speed    float64
}

func (e *LarsonEffect) Render(s *strip.LEDStrip) {
	if e.Speed != 0 {
		e.Position = round(OscillateBetween(0, float64(s.Size-2), e.Speed))
	}

	color := e.getColor()
	s.SetPixel(e.Position, color)
	s.SetPixel(e.Position+1, color)
}

func (e *LarsonEffect) getColor() []uint8 {
	if len(e.Color) == 0 {
		return []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	}

	return e.Color
}
