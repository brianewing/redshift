package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Layer struct {
	Size    int
	Offset  int
	Effects []Effect

	virtualStrip *strip.LEDStrip

	Blend Blend
}

func (e *Layer) Render(s *strip.LEDStrip) {
	if e.virtualStrip == nil {
		if e.Size == 0 {
			e.Size = s.Size
		}

		e.virtualStrip = strip.New(e.Size)
		e.Blend.Buffer = e.virtualStrip.Buffer
	}

	for _, effect := range e.Effects {
		effect.Render(e.virtualStrip)
	}

	e.Blend.Render(s)
}
