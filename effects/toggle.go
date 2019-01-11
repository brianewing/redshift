package effects

import "github.com/brianewing/redshift/strip"

type Toggle struct {
	Effects EffectSet
	Enabled bool
}

func NewToggle() *Toggle {
	return &Toggle{Enabled: true}
}

func (e *Toggle) Init() {
	e.Effects.Init()
}

func (e *Toggle) Destroy() {
	e.Effects.Destroy()
}

func (e *Toggle) Render(s *strip.LEDStrip) {
	if e.Enabled {
		e.Effects.Render(s)
	}
}
