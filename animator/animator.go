package animator

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
)

type Animator struct {
	Strip       *strip.LEDStrip
	Effects     effects.EffectSet
	PostEffects effects.EffectSet

	Running bool
	didInit bool

	Performance *Performance
}

func (a *Animator) Init() {
	a.Performance = NewPerformance()
}

func (a *Animator) Run(interval time.Duration) {
	a.Running = true
	log.Println("Running")

	// TODO: ensure race safety here.
	//       a.Running may be set to false,
	//       then true in another Go routine
	//       before the for loop here can check it.
	ticker := time.NewTicker(interval)
	for ; a.Running; <-ticker.C {
		a.Strip.Lock()
		a.Render()
		a.Strip.Unlock()
	}
	ticker.Stop()
}

func (a *Animator) Render() {
	if !a.didInit {
		a.Effects.InitWithStrip(a.Strip)
		a.PostEffects.InitWithStrip(a.Strip)
		a.didInit = true
	}

	t := time.Now()

	a.Effects.Render(a.Strip)
	a.PostEffects.Render(a.Strip)

	if a.Performance != nil {
		a.Performance.Tick(time.Now().Sub(t))
	}
}

func (a *Animator) SetEffects(newEffects effects.EffectSet) {
	newEffects.InitWithStrip(a.Strip)
	newEffects.Render(strip.New(a.Strip.Size))

	a.Strip.Lock()
	defer a.Strip.Unlock()

	a.Effects.Destroy()
	a.Effects = newEffects

	if newEffects != nil {
		y, _ := effects.MarshalYAML(newEffects)
		ioutil.WriteFile("effects.yaml", y, 0644)
	}
}

func (a *Animator) Finish() {
	a.SetEffects(nil)
	a.Running = false
}
