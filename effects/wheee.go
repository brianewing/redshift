package effects

import "github.com/brianewing/redshift/strip"

type Wheee struct {
	N, i, x int
	Blend   Blend
}

func NewWheee() *Wheee {
	return &Wheee{N: 50, Blend: *NewBlend()}
}

func (e *Wheee) Render(s *strip.LEDStrip) {
	tmp := make(strip.Buffer, len(s.Buffer))

	for i, led := range s.Buffer {
		tmp[i] = make(strip.LED, len(led))
		copy(tmp[i], led)
	}

	e.Blend.Buffer = tmp.Rotate(e.i%len(s.Buffer), false)
	e.Blend.Render(s)

	if e.x++; e.N > 0 && e.x%e.N == 0 {
		e.i++
	}
}