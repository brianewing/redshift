package effects

import (
	"log"

	"github.com/brianewing/redshift/strip"
)

type Resettable struct {
	Effects EffectSet

	Trigger, lastValue int
}

func (e *Resettable) Init(s *strip.LEDStrip) { e.Effects.InitWithStrip(s) }
func (e *Resettable) Destroy()               { e.Effects.Destroy() }

func (e *Resettable) Render(s *strip.LEDStrip) {
	if e.Trigger != e.lastValue {
		e.reset(s)
	}
	e.lastValue = e.Trigger

	e.Effects.Render(s)
}

func (e *Resettable) reset(s *strip.LEDStrip) {
	tmpJSON, _ := MarshalJSON(e.Effects)
	newEffects, _ := UnmarshalJSON(tmpJSON)

	log.Println("reset...", newEffects)
	newEffects.InitWithStrip(s)

	e.Effects.Destroy()
	e.Effects = newEffects
}
