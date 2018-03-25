package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Mirror struct {
	Effects EffectSet
	Size    int

	BlendA Blend
	BlendB Blend

	virtualStrip *strip.LEDStrip
}

func NewMirror() *Mirror {
	return &Mirror{
		BlendA: Blend{
			Reverse: true,
		},
	}
}

func (e *Mirror) Render(s *strip.LEDStrip) {
	if e.virtualStrip == nil {
		if e.Size == 0 {
			e.Size = s.Size
		}

		e.virtualStrip = strip.New(e.Size / 2)

		e.BlendA.Buffer = e.virtualStrip.Buffer
		e.BlendB.Buffer = e.virtualStrip.Buffer
		e.BlendB.Offset = len(e.virtualStrip.Buffer)
	}

	for _, effect := range e.Effects {
		effect.Render(e.virtualStrip)
	}

	e.BlendA.Render(s)
	e.BlendB.Render(s)
}
