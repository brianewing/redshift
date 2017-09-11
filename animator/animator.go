package animator

import (
	"time"
	"redshift/strip"
	"redshift/effects"
)

type Animator struct {
	Strip *strip.LEDStrip
	Effects []effects.Effect
	Interval time.Duration

	Running bool
}

func (a *Animator) Run() {
	a.Running = true
	mutex := &a.Strip.Sync

	for a.Running {
		mutex.Lock()
		a.render()
		mutex.Unlock()

		time.Sleep(a.Interval)
	}
}

func (a *Animator) render() {
	for _, effect := range a.Effects {
		effect.Render(a.Strip)
	}
}

func (a *Animator) Render() {
	a.render()
}

