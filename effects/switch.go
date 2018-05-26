package effects

import "github.com/brianewing/redshift/strip"

type Switch struct {
	Effects   EffectSet
	Selection int
}

func (e *Switch) Init() {
	e.Effects.Init()
}

func (e *Switch) Render(s *strip.LEDStrip) {
	if e.Selection < len(e.Effects) {
		e.Effects[e.Selection].Render(s)
	}
}
