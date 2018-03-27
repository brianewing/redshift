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
	a.init.Do(a.Effects.InitAll)

	a.Effects.Render(a.Strip)
	a.PostEffects.Render(a.Strip)
}

func (a *Animator) SetEffects(newEffects effects.EffectSet) {
	a.Effects.DestroyAll()
	a.Effects = newEffects
	a.init = sync.Once{}
}

func (a *Animator) GetEffects() effects.EffectSet {
	return a.Effects
}
