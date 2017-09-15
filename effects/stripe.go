package effects

import "redshift/strip"

type Stripe struct {
	Color []int
	N int
}

func (e *Stripe) Render(s *strip.LEDStrip) {
	if e.N == 0 { e.N = 1 }

	if len(e.Color) == 0 {
		e.Color = []int{255,255,255}
	}

	for i, led := range s.Buffer {
		if i % e.N == 0 {
			copy(led, e.Color)
		}
	}
}
