package effects

import "github.com/brianewing/redshift/strip"

type Strobe struct {
	N, i    uint
	Reverse bool
}

func NewStrobe() *Strobe {
	return &Strobe{N: 5, Reverse: false}
}

func (e *Strobe) Render(s *strip.LEDStrip) {
	if e.shouldBeOff() {
		s.Clear()
	}
	e.i++
}

func (e *Strobe) shouldBeOff() bool {
	if e.N == 0 {
		return false
	} else if e.Reverse {
		return e.i%e.N == 0
	} else {
		return e.i%e.N != 0
	}
}