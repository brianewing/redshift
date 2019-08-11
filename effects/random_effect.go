package effects

import (
	"github.com/brianewing/redshift/strip"
)

type RandomEffect struct {
	i, N  int
	color []uint8
}

func (e *RandomEffect) Init() {
	if e.N == 0 {
		e.N = 10
	}
	e.color = (&Mood{}).newColor()
}

func (e *RandomEffect) Render(s *strip.LEDStrip) {
	e.i++
	if e.i%e.N == 0 {
		e.color = (&Mood{}).newColor()
		// s.Randomize()
	}
	(&Fill{Color: e.color}).Render(s)
}
