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
	case "Buffer":
		return &Buffer{}
	case "Clear":
		return &Clear{}
	case "External":
		return &External{}
	case "Fill":
		return &Fill{}
	case "Layer":
		return &Layer{}
	case "LarsonEffect":
		return &LarsonEffect{}
	case "MoodEffect":
		return &MoodEffect{}
	case "RaceTestEffect":
		return &RaceTestEffect{}
	case "RainbowEffect":
		return &RainbowEffect{}
	case "RandomEffect":
		return &RandomEffect{}
	case "Stripe":
		return &Stripe{}
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
