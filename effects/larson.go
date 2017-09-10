package effects

import (
	"math/rand"
	"redshift/strip"
)

type LarsonEffect struct {
	Color []int

	velocity int
	position int
}

func (e *LarsonEffect) Render(s *strip.LEDStrip) {
	if e.position <= 1 {
		e.velocity = 1
	} else if e.position == s.Size - 1 {
		e.velocity = -1
	}

	s.Buffer[e.position] = e.getColor()

	if e.position != 0 {
		s.Buffer[e.position-1] = s.Buffer[e.position]
	} else {
		s.Buffer[e.position+1] = s.Buffer[e.position]
	}

	e.position += e.velocity
}

func (e *LarsonEffect) getColor() []int {
	if len(e.Color) == 0 {
		return []int{rand.Intn(255), rand.Intn(255), rand.Intn(255)}
	}

	return e.Color
}
