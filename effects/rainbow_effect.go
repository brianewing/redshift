package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/lucasb-eyer/go-colorful"
	"math/rand"
)

// todo: imagine wheel was huge, say 1000-5000 steps..
// ^^ Rotate it and sample e.g. each (wSize/sSize)th step to cover all hues
// ^^ This would allow for finer speed controls (e.g. midi knobs) and richer color changes
// ^^ Size of the wheel would determine the range/granularity of the speed control..

type RainbowEffect struct {
	Reverse bool
	Size    uint
	Speed   float64
	Dynamic bool

	wheel   [][]uint8
	halting bool
}

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	if e.wheel == nil {
		size := e.getSize()
		e.wheel = generateWheel(size)
	} else if e.Dynamic {
		e.adjustParameters()
	}

	phase := int(CycleBetween(0, float64(len(e.wheel)), e.Speed))
	(&Buffer{Buffer: rotateBuffer(e.wheel, phase, e.Reverse)}).Render(s)
}

func (e *RainbowEffect) adjustParameters() {
	if x := rand.Intn(300); x == 43 {
		e.Reverse = !e.Reverse
	} else if x == 89 {
		e.Speed += 1
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

func (e *RainbowEffect) getSize() uint {
	if e.Size == 0 {
		e.Size = 150
	}
	return e.Size
}

func generateWheel(size uint) [][]uint8 {
	wheel := make([][]uint8, size)
	for i := range wheel {
		hue := float64(i) / float64(len(wheel)) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = []uint8{r, g, b}
	}
	return wheel
}
