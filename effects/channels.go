package effects

import "github.com/brianewing/redshift/strip"

type Channels struct {
	R uint8
	G uint8
	B uint8
}

func NewChannels() *Channels {
	return &Channels{R: 0, G: 1, B: 2}
}

func (e *Channels) Render(strip *strip.LEDStrip) {
	if e.R > 2 || e.G > 2 || e.B > 2 {
		return
	}

	for _, led := range strip.Buffer {
		r, g, b := led[0], led[1], led[2]

		led[0], led[1], led[2] = 0, 0, 0

		led[e.R] = r
		led[e.G] = g
		led[e.B] = b
	}
}
