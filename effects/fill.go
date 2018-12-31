package effects

import "github.com/brianewing/redshift/strip"

type Fill struct {
	Color strip.LED
}

func (e *Fill) Render(s *strip.LEDStrip) {
	for _, led := range s.Buffer {
		copy(led, e.Color)
	}
}
