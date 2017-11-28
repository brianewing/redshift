package animator

import (
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"time"
)

type Animator struct {
	Strip       *strip.LEDStrip
	Effects     []effects.Effect
	PostEffects []effects.Effect

	Running bool
}

func (a *Animator) Run(interval time.Duration) {
	a.Running = true
	mutex := a.Strip

	ticker := time.NewTicker(interval)
	for a.Running {
		mutex.Lock()
		a.Render()
		mutex.Unlock()
		<-ticker.C
	}
}

func (a *Animator) Render() {
	for _, effect := range a.Effects {
		effect.Render(a.Strip)
	}
	for _, postEffect := range a.PostEffects {
		postEffect.Render(a.Strip)
	}
}
