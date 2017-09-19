package effects

import (
	"redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type Buffer struct {
	Buffer [][]uint8
}

func (e *Buffer) Render(strip *strip.LEDStrip) {
	for i, led := range strip.Buffer {
		if i == len(e.Buffer) {
			break
		} else if led[0] == 0 && led[1] == 0 && led[2] == 0  {
			copy(strip.Buffer[i], e.Buffer[i])
		} else {
			copy(strip.Buffer[i], blendRgb(strip.Buffer[i], e.Buffer[i]))
		}
	}
}

func blendRgb(c1 []uint8, c2 []uint8) []uint8 {
	c1_, c2_ := colorfulRgb(c1), colorfulRgb(c2)
	r, g, b := c1_.BlendRgb(c2_, 0.5).Clamped().RGB255()
	return []uint8{r, g, b}
}

func blendHcl(c1 []uint8, c2 []uint8) []uint8 {
	c1_, c2_ := colorfulRgb(c1), colorfulRgb(c2)
	r, g, b := c1_.BlendHcl(c2_, 0.75).Clamped().RGB255()
	return []uint8{r, g, b}
}

func colorfulRgb(c []uint8) colorful.Color {
	r, g, b := float64(c[0]), float64(c[1]), float64(c[2])
	return colorful.Color{R: r / 255.0, G: g / 255.0, B: b / 255.0}
}

