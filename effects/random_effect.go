package effects

import (
	"redshift/strip"
)

type RandomEffect struct {}

func (e *RandomEffect) Render(s *strip.LEDStrip) {
	s.Randomize()
}
