package effects

import "redshift/strip"

type Fill struct {
	Color []uint8
}

func (e *Fill) Render(s *strip.LEDStrip) {
	for _, led := range s.Buffer {
		copy(led, e.Color)
	}
}
