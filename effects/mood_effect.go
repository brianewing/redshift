package effects

import (
	"redshift/strip"
	"math/rand"
)

type MoodEffect struct {
	fillEffect *Fill
	brightnessEffect *Brightness
	layer *Layer
}

func (e *MoodEffect) Render(s *strip.LEDStrip) {
	if e.layer == nil {
		e.fillEffect = &Fill{}
		e.brightnessEffect = &Brightness{}
		e.layer = &Layer{
			Effects: []Effect{
				e.fillEffect,
				e.brightnessEffect,
			},
		}
	}

	if e.brightnessEffect.Brightness <= 1 {
		e.fillEffect.Color = e.newColor()
	}

	e.layer.Render(s)
}

func (e *MoodEffect) newColor() []uint8 {
	return []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
}