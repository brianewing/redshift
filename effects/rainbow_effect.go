package effects

import (
	"math"

	"github.com/brianewing/redshift/strip"
	colorful "github.com/lucasb-eyer/go-colorful"
)

type RainbowEffect struct {
	Size    uint16
	Depth   uint8   `max:20"`
	Speed   float64 `max:"10"`
	Blend   *Blend  `advanced`
	Reverse bool

	NewMethod bool

	wheel strip.Buffer
}

func NewRainbowEffect() *RainbowEffect {
	return &RainbowEffect{
		Size:      0,
		Depth:     5,
		Speed:     0.1,
		Blend:     NewBlend(),
		NewMethod: true,
	}
}

func (e *RainbowEffect) Render(s *strip.LEDStrip) {
	if !e.NewMethod {
		e.oldRender(s)
		return
	}

	if e.Size == 0 {
		e.Size = uint16(s.Size)
	}

	if len(e.Blend.Buffer) != int(e.Size) {
		e.Blend.Buffer = strip.NewBuffer(int(e.Size))
	}

	// phase := float64(round(CycleBetween(0, 360.0, e.Speed)))
	speed := e.Speed

	var phase float64
	if e.Reverse {
		phase = CycleBetween(0, 360.0, speed)
	} else {
		phase = CycleBetween(360, 0, speed)
	}

	for i, led := range e.Blend.Buffer {
		color := colorful.Hsv(
			// phase,
			math.Mod(phase+float64(i)/float64(e.Size)*360, 360),
			// math.Abs(math.Mod(phase-float64(i)/float64(e.Size)*360, 360)),
			1,
			1)
		led[0], led[1], led[2] = color.Clamped().RGB255()
	}

	e.Blend.Render(s)
}

func (e *RainbowEffect) oldRender(s *strip.LEDStrip) {
	if e.Depth == 0 {
		e.Depth = 1
	} else if e.Depth > 20 {
		e.Depth = 20
	}

	if e.Size == 0 {
		e.Size = uint16(s.Size)
	}

	steps := e.Size * uint16(e.Depth)

	if e.wheel == nil || len(e.wheel) != int(steps) {
		e.wheel = strip.Buffer(generateWheel(steps))
	}

	phase := round(CycleBetween(0, float64(len(e.wheel)), e.Speed))

	e.Blend.Reverse = e.Reverse
	e.Blend.Buffer = e.wheel.Rotate(phase, false).Sample(int(e.Size))
	e.Blend.Render(s)
}

func generateWheel(size uint16) strip.Buffer {
	wheel := make(strip.Buffer, size)

	for i := range wheel {
		hue := float64(i) / float64(size) * 360
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		wheel[i] = strip.LED{r, g, b}
	}

	return wheel
}
