package effects

import (
	"math/rand"
	"github.com/lucasb-eyer/go-colorful"
	"redshift/strip"
)

type RainbowEffect struct {
	Static bool
	Reverse bool

	wheel [][]int
}

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	if e.wheel == nil {
		e.wheel = generateWheel(s)
	} else if !e.Static {
		e.wheel = rotateWheel(e.wheel, e.Reverse)
	}

	if rand.Intn(200) == 43 {
		e.Reverse = !e.Reverse
	}

	(&Buffer{Buffer: e.wheel}).Render(s)
}

func generateWheel(s *strip.LEDStrip) [][]int {
	wheel := make([][]int, 150)
	for i := range wheel {
		hue := float64(i) / float64(len(wheel)) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = []int{int(r), int(g), int(b)}
	}
	return wheel
}

func rotateWheel(wheel [][]int, reverse bool) [][]int {
	if reverse {
		head, wheel := wheel[0], wheel[1:]
		return append(wheel, head)
	} else {
		wheel, tail := wheel[:len(wheel)-1], wheel[len(wheel)-1]
		return append([][]int{tail}, wheel...)
	}
}
