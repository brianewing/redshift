package effects

import (
	"redshift/strip"
)

type Effect interface {
	Render(strip *strip.LEDStrip)
}

func newEffectByName(name string) Effect {
	switch name {
		case "BlueEffect": return &BlueEffect{}
		case "Brightness": return &Brightness{}
		case "Buffer": return &Buffer{}
		case "Clear": return &Clear{}
		case "Combine": return &Combine{}
		case "LarsonEffect": return &LarsonEffect{}
		case "RaceTestEffect": return &RaceTestEffect{}
		case "RainbowEffect": return &RainbowEffect{}
		case "RandomEffect": return &RandomEffect{}
		case "RotateBuffer": return &RotateBuffer{}
		case "Stripe": return &Stripe{}
		default: return nil
	}
}
