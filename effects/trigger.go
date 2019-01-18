package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Trigger struct {
	Effects EffectSet
	Value   int

	lastSeenValue int
}

func (t *Trigger) Init() {
	t.Effects.Init()
}

func (t *Trigger) Destroy() {
	t.Effects.Destroy()
}

func (t *Trigger) Render(s *strip.LEDStrip) {
	if t.Value != t.lastSeenValue {
		t.Effects.Render(s)
	}

	t.lastSeenValue = t.Value
}
