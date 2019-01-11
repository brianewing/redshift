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

func (a *Animator) Run(interval time.Duration) {
	log.Println("Running")
	a.Running = true
	go a.logFrameRate()

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
		a.Effects.Init()
		a.PostEffects.Init()
		a.didInit = true
	}

	t1 := time.Now()
	a.Effects.Render(a.Strip)
	a.PostEffects.Render(a.Strip)
	if a.Performance != nil {
		a.Performance.Tick(time.Now().Sub(t1))
	}
}

func (a *Animator) SetEffects(newEffects effects.EffectSet) {
	newEffects.Init()
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

func (a *Animator) logFrameRate() {
	a.Performance = NewPerformance()

	ticker := time.NewTicker(1 * time.Second)
	<-ticker.C // wait for first second to pass

	for ; a.Running; <-ticker.C {
		log.Println(a.Performance.String())
		a.Performance.Reset()
	}
	ticker.Stop()
}
