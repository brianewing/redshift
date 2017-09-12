package effects

import "redshift/strip"

type BlueEffect struct {
	value int
	direction int
}

func (e *BlueEffect) Render(s *strip.LEDStrip) {
	if e.value == 255 {
		e.direction = -1
	} else if e.value < 1  {
		e.direction = 1
	}

	e.value += e.direction

	for _, led := range s.Buffer {
		led[2] = e.value
	}
}