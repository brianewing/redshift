package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Effect interface {
	Render(strip *strip.LEDStrip)
}

type Initable interface { Init() }
type Destroyable interface { Destroy() }

func NewByName(name string) Effect {
	switch name {
	case "BlueEffect":
		return &BlueEffect{}
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
		"MoodEffect",
		"Null",
		"RainbowEffect",
		"RandomEffect",
		"Stripe",
	}
}
