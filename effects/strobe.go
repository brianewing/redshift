package effects

import "github.com/brianewing/redshift/strip"

type Strobe struct {
	N, i int
}

func NewStrobe() *Strobe {
	return &Strobe{N: 10}
}

func (e *Strobe) Render(s *strip.LEDStrip) {
	if e.N > 0 && e.i % e.N != 0 {
		s.Clear()
	}

	e.i++
}
