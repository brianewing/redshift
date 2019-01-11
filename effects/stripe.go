package effects

import "github.com/brianewing/redshift/strip"

type Stripe struct {
	Color strip.LED
	N     int
}

func NewStripe() *Stripe {
	return &Stripe{
		Color: strip.LED{255, 255, 255},
		N:     2,
	}
}

func (e *Stripe) Render(s *strip.LEDStrip) {
	for i, led := range s.Buffer {
		if e.N == 0 || i%e.N == 0 {
			copy(led, e.Color)
		}
	}
}
