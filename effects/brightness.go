package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type Brightness struct {
	Level int // 0-255, todo: return to uint8
}

func (e *Brightness) Render(s *strip.LEDStrip) {
	for _, color := range s.Buffer {
		//applyBasic(color, uint8(e.Brightness))
		applyHsl(color, uint8(e.Level))
	}
}

func applyHsl(color []uint8, brightness uint8) {
	r, g, b := float64(color[0]), float64(color[1]), float64(color[2])
	c := colorful.Color{R: r / 255.0, G: g / 255.0, B: b / 255.0}

	h, s, v := c.Hsl()
	v *= (float64(brightness) / 255.0)

	newR, newG, newB := colorful.Hsl(h, s, v).Clamped().RGB255()
	copy(color, []uint8{newR, newG, newB})
}

// this doesn't treat colours evenly but may be useful
func applyBasic(color []uint8, brightness uint8) {
	for i, v := range color {
		if v > brightness {
			color[i] = brightness
		}
	}
}
