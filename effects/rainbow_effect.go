package effects

import (
	"math/rand"
	"github.com/lucasb-eyer/go-colorful"
	"redshift/strip"
)

type RainbowEffect struct {
	Reverse bool
	Size int
	Speed int
	Dynamic bool

	wheel [][]int
	halting bool
}

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	if e.wheel == nil {
		size := e.getSize()
		e.wheel = generateWheel(size)
	} else if e.Dynamic {
		e.adjustParameters()
	}

	(&RotateBuffer{Buffer: &e.wheel, Count: e.Speed, Reverse: e.Reverse}).Render(s)
}

func (e *RainbowEffect) adjustParameters() {
	if x := rand.Intn(200); x == 43 {
		e.Reverse = !e.Reverse
	} else if x == 89 {
		e.Speed += 2
	} else if x == 48 && e.Speed > 1 {
		e.Speed -= 1
	} else if x == 31 && e.Speed > 1 {
		e.halting = true
	} else if e.halting && e.Speed > 1 && x > 185 {
		e.Speed -= 1
	} else if e.halting && e.Speed == 1 {
		e.halting = false
	}
}

func (e *RainbowEffect) getSize() int {
	if e.Size == 0 {
		e.Size = 150
	}
	return e.Size
}

func generateWheel(size int) [][]int {
	wheel := make([][]int, size)
	for i := range wheel {
		hue := float64(i) / float64(len(wheel)) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = []int{int(r), int(g), int(b)}
	}
	return wheel
}
