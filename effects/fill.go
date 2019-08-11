package effects

import "github.com/brianewing/redshift/strip"

type Fill struct {
	Color strip.LED
}

func NewFill() *Fill {
	return &Fill{
		Color: strip.LED{0, 0, 0},
	}
}

func (e *Fill) Render(s *strip.LEDStrip) {
	if e.Color[0] == 0 && e.Color[1] == 0 && e.Color[2] == 0 {
		return
	}

	for _, led := range s.Buffer {
		copy(led, e.Color)
	}
}
