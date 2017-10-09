package effects

import (
	"redshift/strip"
)

type Layer struct {
	Effects []Effect

	virtualStrip *strip.LEDStrip
}

func (e *Layer) Render(s *strip.LEDStrip) {
	if e.virtualStrip == nil {
		e.virtualStrip = strip.New(s.Size)
	}

	for _, effect := range e.Effects {
		effect.Render(e.virtualStrip)
	}

	(&Buffer{Buffer: e.virtualStrip.Buffer}).Render(s)
}
