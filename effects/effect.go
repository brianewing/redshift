package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Effect interface {
	Render(strip *strip.LEDStrip)
}

type Initable interface { Init() }
type Destroyable interface { Destroy() }

type EffectEnvelope struct {
	Effect
	Controls []struct{}
}

type EffectSet []EffectEnvelope

func (s EffectSet) Render(strip *strip.LEDStrip) {
	for _, effect := range s {
		effect.Render(strip)
	}
}

// calls Init() on initable effects
func (s EffectSet) InitAll() {
	for _, envelope := range s {
		if initable, ok := envelope.Effect.(Initable); ok {
			initable.Init()
		}
	}
}

// calls Destroy() on destroyable effects
func (s EffectSet) DestroyAll() {
	for _, envelope := range s {
		if destroyable, ok := envelope.Effect.(Destroyable); ok {
			destroyable.Destroy()
		}
	}
}

func NewByName(name string) Effect {
	switch name {
	case "BlueEffect":
		return NewBlueEffect()
	case "Brightness":
		return &Brightness{}
	case "Blend":
		return &Blend{}
	case "Clear":
		return &Clear{}
	case "External":
		return &External{}
	case "Fill":
		return &Fill{}
	case "Layer":
		return &Layer{}
	case "LarsonEffect":
		return NewLarsonEffect()
	case "Mirror":
		return &Mirror{}
	case "MoodEffect":
		return &MoodEffect{}
	case "RaceTestEffect":
		return &RaceTestEffect{}
	case "RainbowEffect":
		return NewRainbowEffect()
	case "RandomEffect":
		return &RandomEffect{}
	case "Stripe":
		return NewStripe()
	case "Strobe":
		return NewStrobe()
	case "Switch":
		return &Switch{}
	case "Toggle":
		return &Toggle{}
	default:
		return &Null{}
	}
}

func Names() []string {
	return []string{
		"BlueEffect",
		"Brightness",
		"Buffer",
		"Clear",
		"External",
		"Fill",
		"Layer",
		"LarsonEffect",
		"Mirror",
		"MoodEffect",
		"Null",
		"RainbowEffect",
		"RandomEffect",
		"Stripe",
		"Strobe",
		"Switch",
		"Toggle",
	}
}
