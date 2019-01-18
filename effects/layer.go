package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Layer struct {
	Blend   Blend
	Effects EffectSet
	Size    int
	Offset  int

	virtualStrip *strip.LEDStrip
}

func NewLayer() *Layer {
	return &Layer{
		Blend:   *NewBlend(),
		Effects: EffectSet{EffectEnvelope{Effect: &Clear{}}},
	}
}

func (e *Layer) Init() {
	e.Effects.Init()
}

func (e *Layer) Destroy() {
	e.Effects.Destroy()
}

func (e *Layer) Render(s *strip.LEDStrip) {
	if e.virtualStrip == nil || e.virtualStrip.Size != e.Size {
		if e.Size == 0 {
			e.Size = s.Size
		}

		e.virtualStrip = strip.New(e.Size)
		e.Blend.Buffer = e.virtualStrip.Buffer
	}

	e.Blend.Offset = e.Offset

	e.Effects.Render(e.virtualStrip)
	e.Blend.Render(s)
}
