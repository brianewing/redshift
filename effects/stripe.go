package effects

import "github.com/brianewing/redshift/strip"

type Stripe struct {
	Color []uint8
	N     int
}

func (e *Stripe) Render(s *strip.LEDStrip) {
	if e.N == 0 {
		e.N = 2
	}

	if len(e.Color) == 0 {
		e.Color = []uint8{255, 255, 255}
	}

	for i, led := range s.Buffer {
		if i%e.N == 0 {
			copy(led, e.Color)
		}
	}
}
