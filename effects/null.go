package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Null struct{}

func (e *Null) Render(strip *strip.LEDStrip) {}
