package effects

import (
	"github.com/brianewing/redshift/strip"
	"math/rand"
)

type LarsonEffect struct {
	Color    []uint8
	Position int
	velocity int
}

func (e *LarsonEffect) Render(s *strip.LEDStrip) {
	if e.Position <= 1 || e.velocity == 0 {
		e.velocity = 1
	} else if e.Position >= s.Size-1 {
		e.velocity = -1
	}

	color := e.getColor()
	s.SetPixel(e.Position, color)

	if e.Position != 0 {
		s.SetPixel(e.Position-1, color)
	} else {
		s.SetPixel(e.Position+1, color)
	}

	e.Position += e.velocity
}

func (e *LarsonEffect) getColor() []uint8 {
	if len(e.Color) == 0 {
		return []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	}

	return e.Color
}
