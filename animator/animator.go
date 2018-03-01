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
	didInit bool
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
	ticker.Stop()
}

func (a *Animator) Render() {
	if !a.didInit {
		initEffects(a.Effects)
		a.didInit = true
	}
	for _, effect := range a.Effects {
		effect.Render(a.Strip)
	}
	for _, postEffect := range a.PostEffects {
		postEffect.Render(a.Strip)
	}
}

func (a *Animator) SetEffects(newEffects []effects.Effect) {
	destroyEffects(a.Effects)
	a.Effects = newEffects
	a.didInit = false
}

func (a *Animator) GetEffects() []effects.Effect {
	return a.Effects
}

// calls Init() on initable effects
func initEffects(effects_ []effects.Effect) {
	for _, e := range effects_ {
		if initable, ok := e.(effects.Initable); ok {
			initable.Init()
		}
	}
}

// calls Destroy() on destroyable effects
func destroyEffects(effects_ []effects.Effect) {
	for _, e := range effects_ {
		if destroyable, ok := e.(effects.Destroyable); ok {
			destroyable.Destroy()
		}
	}
}
