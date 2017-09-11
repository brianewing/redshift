package animator

import (
	"time"
	"redshift/strip"
	"redshift/effects"
)

type Animator struct {
	Strip *strip.LEDStrip
	Effects []effects.Effect

	Running bool
}

func (a *Animator) Run(interval time.Duration) {
	a.Running = true
	mutex := &a.Strip.Sync

	for a.Running {
		mutex.Lock()
		a.render()
		mutex.Unlock()

		time.Sleep(interval)
	}
}

func (a *Animator) render() {
	for _, effect := range a.Effects {
		effect.Render(a.Strip)
	}
}
