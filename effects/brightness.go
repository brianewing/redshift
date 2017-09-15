package effects

import (
	"redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
)

type Brightness struct {
	Brightness int
	velocity int
}

// Acts as a periodic fade when Brightness is 0 or 255
// Eventually this will be controlled externally..
// e.g. by an animation function, gui or midi controller

func (e *Brightness) Render(s *strip.LEDStrip) {
	if e.Brightness > 254 {
		e.Brightness = 255
		e.velocity = -2
	} else if e.Brightness < 1 {
		e.Brightness = 0
		e.velocity = 2
	}

	for _, color := range s.Buffer {
		//applyBasic(color, int(uint8(e.Brightness)))
		applyHsv(color, e.Brightness)
	}

	e.Brightness += e.velocity
}

func applyHsv(color []int, brightness int) {
	r, g, b := float64(color[0]), float64(color[1]), float64(color[2])
	c := colorful.Color{R: r / 255.0, G: g / 255.0, B: b / 255.0}

	h, s, v := c.Hsl()
	v *= (float64(brightness) / 255.0)

	newR, newG, newB := colorful.Hsl(h, s, v).Clamped().RGB255()
	copy(color, []int{int(newR), int(newG), int(newB)})
}

// this doesn't treat colours evenly but may be useful
func applyBasic(color []int, brightness int) {
	for i, v := range color {
		if v > brightness {
			color[i] = brightness
		}
	}
}
