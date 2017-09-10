package effects

import "redshift/strip"

type Effect interface {
	Render(strip *strip.LEDStrip)
}

type Clear struct {}

func (e *Clear) Render(strip *strip.LEDStrip) {
	strip.Clear()
}