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

	color := e.getColor()
	s.Buffer[e.position] = color
	//s.Buffer[e.position][0] = color[0]
	//s.Buffer[e.position][1] = color[1]
	//s.Buffer[e.position][2] = color[2]

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

	//return e.Color
	return copyColor(e.Color)
}

func copyColor(c []int) []int {
	c2 := make([]int, len(c))
	copy(c2, c)

	return c2
}
