package animator

import (
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"io/ioutil"
	"log"
	"time"
)

type Animator struct {
	Strip       *strip.LEDStrip
	Effects     effects.EffectSet
	PostEffects effects.EffectSet

	Running bool
	didInit bool

	framesLastPrint time.Time
	framesCount     int

	frameDurations []time.Duration
}

func min(durations []time.Duration) (m time.Duration) {
	if len(durations) > 0 {
		m = durations[0]
		for _, d := range durations {
			if d < m {
				m = d
			}
		}
	}
	return
}

func max(durations []time.Duration) (m time.Duration) {
	if len(durations) > 0 {
		m = durations[0]
		for _, d := range durations {
			if d > m {
				m = d
			}
		}
	}
	return
}

func (a *Animator) logFrameRate() {
	ticker := time.NewTicker(1 * time.Second)

	for ; ; <-ticker.C {
		if a.framesLastPrint.IsZero() {
			a.framesLastPrint = time.Now()
			continue
		}

		if !a.Running {
			return
		}

		a.Strip.Lock()

		framesPerSecond := float64(a.framesCount) / float64((time.Now().Sub(a.framesLastPrint))) * float64(time.Second)
		a.framesCount = 0

		minT := min(a.frameDurations)
		maxT := max(a.frameDurations)

		a.frameDurations = []time.Duration{}
		
		a.Strip.Unlock()

		log.Println("frames per second", framesPerSecond, "min", minT, "max", maxT)

		a.framesLastPrint = time.Now()
	}

	ticker.Stop()
}

func (a *Animator) Run(interval time.Duration) {
	// go a.logFrameRate()
	a.frameDurations = []time.Duration{}
	ticker := time.NewTicker(interval)
	for a.Running = true; a.Running; <-ticker.C {
		a.Strip.Lock()
		a.framesCount += 1
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
	a.frameDurations = append(a.frameDurations, time.Now().Sub(t1))
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
