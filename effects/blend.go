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

type blendFunction func(a, b []uint8, factor float64) (c []uint8)

func NewBlend() *Blend {
	return &Blend{
		Func:   "rgb",
		Factor: 0.5,
	}
}

func NewBlendFromBuffer(buffer [][]uint8) *Blend {
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

		if isOff(dest) && !e.Force {
			copy(dest, source)
		} else if !isOff(source) || e.Force {
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

func blendRgb(a []uint8, b []uint8, factor float64) (c []uint8) {
	cA, cB := colorfulRgb(a), colorfulRgb(b)
	r, g, b_ := cA.BlendRgb(cB, factor).Clamped().RGB255()
	return []uint8{r, g, b_}
}

func blendHcl(a []uint8, b []uint8, factor float64) (c []uint8) {
	cA, cB := colorfulRgb(a), colorfulRgb(b)
	r, g, b_ := cA.BlendHcl(cB, factor).Clamped().RGB255()
	return []uint8{r, g, b_}
}

func blendLab(a []uint8, b []uint8, factor float64) (c []uint8) {
	cA, cB := colorfulRgb(a), colorfulRgb(b)
	r, g, b_ := cA.BlendLab(cB, factor).Clamped().RGB255()
	return []uint8{r, g, b_}
}

func blendNone(_ []uint8, b []uint8, factor float64) []uint8 {
	return b
}

// Misc buffer functions

func colorfulRgb(c []uint8) colorful.Color {
	r, g, b := float64(c[0]), float64(c[1]), float64(c[2])
	return colorful.Color{R: r / 255.0, G: g / 255.0, B: b / 255.0}
}

func isOff(led []uint8) bool {
	return led[0] == 0 && led[1] == 0 && led[2] == 0
}

// returns a new slice containing the data in buffer rotated by n
func rotateBuffer(buffer [][]uint8, n int, reverse bool) [][]uint8 {
	if reverse {
		head, tail := buffer[0:n], buffer[n:]
		return append(tail, head...)
	} else {
		head, tail := buffer[:len(buffer)-n], buffer[len(buffer)-n:]
		return append(tail, head...)
	}
}

// returns a subset of buffer (n evenly-spaced elements)
func sampleBuffer(buffer [][]uint8, n int) [][]uint8 {
	subset := make([][]uint8, n)
	if n > 0 {
		step := len(buffer) / n
		for i := 0; i < n; i++ {
			subset[i] = buffer[i*step]
		}
	}
	return subset
}
