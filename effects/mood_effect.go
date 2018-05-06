package effects

import (
	"github.com/brianewing/redshift/strip"
	"math/rand"
)

type MoodEffect struct {
	fill       *Fill
	brightness *Brightness
	layer      *Layer

	Speed float64
}

func NewMoodEffect() *MoodEffect {
	return &MoodEffect{
		Speed: 0.1,
	}
}

func (e *MoodEffect) Init() {
	e.fill = &Fill{}
	e.brightness = &Brightness{}
	e.layer = NewLayer()

	e.layer.Effects = EffectSet{
		EffectEnvelope{Effect: e.fill},
		EffectEnvelope{
			Effect: e.brightness,
			Controls: ControlSet{
				ControlEnvelope{
					Control: &TweenControl{
						BaseControl: BaseControl{Field: "Level"},
						Min:   0,
						Max:   255,
						Function: "triangle",
						Speed: e.Speed,
					},
				},
			},
		},
	}

	e.layer.Init()
}

func (e *MoodEffect) Render(s *strip.LEDStrip) {
	if e.brightness.Level <= 5 {
		e.fill.Color = e.newColor()
	}

	e.layer.Render(s)
}

func (e *MoodEffect) newColor() []uint8 {
	return []uint8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
}
