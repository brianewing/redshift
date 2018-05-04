package animator

import (
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"sync"
	"time"
)

type Animator struct {
	Strip       *strip.LEDStrip
	Effects     effects.EffectSet
	PostEffects effects.EffectSet

	Running bool
	init    sync.Once
}

func (a *Animator) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for a.Running = true; a.Running; <-ticker.C {
		a.Strip.Lock()
		a.Render()
		a.Strip.Unlock()
	}
	ticker.Stop()
}

func (a *Animator) Render() {
	a.init.Do(a.Effects.Init)

	a.Effects.Render(a.Strip)
	a.PostEffects.Render(a.Strip)
}

func (a *Animator) SetEffects(newEffects effects.EffectSet) {
	a.Effects.Destroy()
	a.Effects = newEffects
	a.init = sync.Once{} // init again on next Render
}

func (a *Animator) GetEffects() effects.EffectSet {
	return a.Effects
}
