package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type Blend struct {
	Buffer  [][]uint8 `json:"-"`
	Offset  int
	Reverse bool
	Force   bool

	Func string // e.g. rgb, hcl, lab
}

func NewBlend() *Blend {
	return &Blend{
		Func: "rgb",
	}
}

func (e *Blend) Render(strip *strip.LEDStrip) {
	for i := e.Offset; i < len(e.Buffer)+e.Offset && i < strip.Size; i++ {
		source := e.Buffer[i-e.Offset]
		dest := strip.Buffer[i]

		if e.Reverse {
			dest = strip.Buffer[len(e.Buffer)-i-1]
		}

		if isOff(dest) && !e.Force {
			copy(dest, source)
		} else if !isOff(source) || e.Force {
			blendFn := e.getFunction()
			copy(dest, blendFn(dest, source))
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
	return noBlend
}

type blendFunction func(c1, c2 []uint8) (c3 []uint8)

func blendRgb(c1 []uint8, c2 []uint8) []uint8 {
	c1_, c2_ := colorfulRgb(c1), colorfulRgb(c2)
	r, g, b := c1_.BlendRgb(c2_, 0.5).Clamped().RGB255()
	return []uint8{r, g, b}
}

func blendHcl(c1 []uint8, c2 []uint8) []uint8 {
	c1_, c2_ := colorfulRgb(c1), colorfulRgb(c2)
	r, g, b := c1_.BlendHcl(c2_, 0.5).Clamped().RGB255()
	return []uint8{r, g, b}
}

func blendLab(c1 []uint8, c2 []uint8) []uint8 {
	c1_, c2_ := colorfulRgb(c1), colorfulRgb(c2)
	r, g, b := c1_.BlendLab(c2_, 0.5).Clamped().RGB255()
	return []uint8{r, g, b}
}

func noBlend(_ []uint8, c2 []uint8) []uint8 {
	return c2
}

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
