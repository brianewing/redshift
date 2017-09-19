package effects

import (
	"math/rand"
	"redshift/strip"
)

type LarsonEffect struct {
	Color []uint8

	velocity int
	position int
}

func (e *LarsonEffect) Render(s *strip.LEDStrip) {
	if e.position <= 1 {
		e.velocity = 1
	} else if e.position == s.Size - 1 {
		e.velocity = -1
	}

	color := e.getColor()
	s.SetPixel(e.position, color)

	if e.position != 0 {
		s.SetPixel(e.position-1, color)
	} else {
		s.SetPixel(e.position+1, color)
	}

	e.position += e.velocity
}

func (e *LarsonEffect) getColor() []uint8 {
	if len(e.Color) == 0 {
		return []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	}

	return e.Color
}
