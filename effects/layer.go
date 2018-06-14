package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Layer struct {
	Size         int
	Effects      EffectSet
	virtualStrip *strip.LEDStrip

	Blend Blend
}

func NewLayer() *Layer {
	return &Layer{
		Blend: *NewBlend(),
	}
}

func (e *Layer) Init() {
	e.Effects.Init()
}

func (e *Layer) Destroy() {
	e.Effects.Destroy()
}

func (e *Layer) Render(s *strip.LEDStrip) {
	if e.virtualStrip == nil {
		if e.Size == 0 {
			e.Size = s.Size
		}

		e.virtualStrip = strip.New(e.Size)
		e.Blend.Buffer = e.virtualStrip.Buffer
	}

	e.Effects.Render(e.virtualStrip)

	e.Blend.Render(s)
}
