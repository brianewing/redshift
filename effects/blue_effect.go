package effects

import "redshift/strip"

type BlueEffect struct {
	value uint8
	backwards bool
}

func (e *BlueEffect) Render(s *strip.LEDStrip) {
	if e.value == 255 {
		e.backwards = true
	} else if e.value < 1  {
		e.backwards = false
	}

	if e.backwards {
		e.value -= 1
	} else {
		e.value += 1
	}

	for _, led := range s.Buffer {
		led[2] = e.value
	}
}