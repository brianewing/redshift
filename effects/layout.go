package effects

import (
	"math"

	"github.com/brianewing/redshift/strip"
)

type Layout struct {
	Type    string
	Effects EffectSet

	lastType string

	layer        *Layer
	virtualStrip *strip.LEDStrip
}

func NewLayout() *Layout {
	return &Layout{
		Type: "grid", // grid, mirror, none
	}
}

func (e *Layout) Init() {
	e.Effects.Init()
}

func (e *Layout) Destroy() {
	e.Effects.Destroy()
}

func (e *Layout) Render(s *strip.LEDStrip) {
	if e.virtualStrip == nil || e.lastType != e.Type {
		if e.Type == "grid" || e.Type == "line" {
			e.virtualStrip = strip.New(int(math.Sqrt(float64(s.Size))))
		} else if e.Type == "mirror" {
			e.virtualStrip = strip.New(s.Size / 2)
		} else { // "none"
			e.virtualStrip = strip.New(s.Size)
		}
	}

	e.lastType = e.Type

	e.Effects.Render(e.virtualStrip)
	out := strip.New(s.Size)

	if e.Type == "grid" {
		columns := e.virtualStrip.Size
		for i := 0; i < s.Size; i += columns {
			for j := 0; j < columns; j++ {
				source := e.virtualStrip.Buffer[int(i/columns)]
				copy(out.Buffer[i+j], source)
			}
		}
	} else if e.Type == "line" {
		columns := e.virtualStrip.Size
		for i := 0; i < s.Size; i += columns {
			for j := 0; j < columns; j++ {
				if j >= len(e.virtualStrip.Buffer) || i+j >= len(out.Buffer) {
					break
				}
				source := e.virtualStrip.Buffer[j]
				copy(out.Buffer[i+j], source)
			}
		}
	} else if e.Type == "mirror" {
		(&Blend{Offset: e.virtualStrip.Size, Buffer: e.virtualStrip.Buffer}).Render(out)
		(&Blend{Offset: 0, Reverse: true, Buffer: e.virtualStrip.Buffer}).Render(out)
	} else {
		out.Buffer = e.virtualStrip.Buffer
	}

	(&Blend{Func: "none", Buffer: out.Buffer}).Render(s)
}
