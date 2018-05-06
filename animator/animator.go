package animator

import (
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"time"
)

type Animator struct {
	Strip       *strip.LEDStrip
	Effects     effects.EffectSet
	PostEffects effects.EffectSet

	Running bool
	didInit bool
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
	if !a.didInit {
		a.Effects.Init()
		a.PostEffects.Init()
		a.didInit = true
	}

	a.Effects.Render(a.Strip)
	a.PostEffects.Render(a.Strip)
}

func (a *Animator) SetEffects(newEffects effects.EffectSet) {
	a.Effects.Destroy()
	a.Effects = newEffects
	a.didInit = false // init again on next Render
}

func (a *Animator) Finish() {
	a.SetEffects(nil)
	a.Running = false
}
