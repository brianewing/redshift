package effects

import "github.com/brianewing/redshift/strip"

type Whoosh struct {
	N, i, x int
	Reverse bool
	Blend   Blend
}

func NewWhoosh() *Whoosh {
	return &Whoosh{N: 50, Blend: *NewBlend()}
}

func (e *Whoosh) Render(s *strip.LEDStrip) {
	tmp := make(strip.Buffer, len(s.Buffer))

	for i, led := range s.Buffer {
		tmp[i] = make(strip.LED, len(led))
		copy(tmp[i], led)
	}

	e.Blend.Buffer = tmp.Rotate(e.i%len(s.Buffer), e.Reverse)
	e.Blend.Render(s)

	if e.x++; e.N > 0 && e.x%e.N == 0 {
		e.i++
	}
}
