package effects

import "github.com/brianewing/redshift/strip"

type Clear struct{}

func (e *Clear) Render(strip *strip.LEDStrip) {
	strip.Clear()
}
