package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type Blend struct {
	Buffer  strip.Buffer `json:"-"`
	Offset  int
	Reverse bool
	Force   bool

	Func   string // e.g. rgb, hcl, lab
	Factor float64
}

type blendFunction func(a, b strip.LED, factor float64) (c strip.LED)

func NewBlend() *Blend {
	return &Blend{
		Func:   "rgb",
		Factor: 0.5,
	}
}

func NewBlendFromBuffer(buffer strip.Buffer) *Blend {
	b := NewBlend()
	b.Buffer = buffer
	return b
}

func (e *Blend) Render(strip *strip.LEDStrip) {
	for i := 0; i < len(e.Buffer) && i+e.Offset < strip.Size; i++ {
		source := e.Buffer[i]
		dest := strip.Buffer[i+e.Offset]

		if e.Reverse {
			dest = strip.Buffer[len(e.Buffer)+e.Offset-i-1]
		}

		if dest.IsOff() && !e.Force {
			copy(dest, source)
		} else if !source.IsOff() || e.Force {
			blendFn := e.getFunction()
			copy(dest, blendFn(dest, source, e.Factor))
		}
	}
}

func (e *Blend) getFunction() blendFunction {
	switch e.Func {
	case "hcl":
		return blendHcl
	case "lab":
		return blendLab
	case "rgb":
		return blendRgb
	}
	return blendNone
}

// Blend functions

func blendRgb(a strip.LED, b strip.LED, factor float64) (c strip.LED) {
	cA, cB := colorfulRgb(a), colorfulRgb(b)
	r, g, b_ := cA.BlendRgb(cB, factor).Clamped().RGB255()
	return strip.LED{r, g, b_}
}

func blendHcl(a strip.LED, b strip.LED, factor float64) (c strip.LED) {
	cA, cB := colorfulRgb(a), colorfulRgb(b)
	r, g, b_ := cA.BlendHcl(cB, factor).Clamped().RGB255()
	return strip.LED{r, g, b_}
}

func blendLab(a strip.LED, b strip.LED, factor float64) (c strip.LED) {
	cA, cB := colorfulRgb(a), colorfulRgb(b)
	r, g, b_ := cA.BlendLab(cB, factor).Clamped().RGB255()
	return strip.LED{r, g, b_}
}

func blendNone(_ strip.LED, b strip.LED, factor float64) strip.LED {
	return b
}

// Misc buffer functions

func colorfulRgb(c strip.LED) colorful.Color {
	r, g, b := float64(c[0]), float64(c[1]), float64(c[2])
	return colorful.Color{R: r / 255.0, G: g / 255.0, B: b / 255.0}
}
