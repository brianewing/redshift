package effects

import "github.com/brianewing/redshift/strip"

type BlueEffect struct {
	Speed float64
}

func (e *BlueEffect) Render(s *strip.LEDStrip) {
	value := uint8(OscillateBetween(0, 255, e.Speed))
	for _, led := range s.Buffer {
		led[2] = value
	}
}
