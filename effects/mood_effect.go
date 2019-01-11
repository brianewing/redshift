package effects

import (
	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/brianewing/redshift/strip"
)

type MoodEffect struct {
	fill       *Fill
	brightness *Brightness
	layer      *Layer

	Speed float64
}

func NewMoodEffect() *MoodEffect {
	return &MoodEffect{
		// Speed: 0.1,
		Speed: 0.3,
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
						Min:         1,
						Max:         255,
						Function:    "sin",
						Speed:       e.Speed,
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
	r, g, b := colorful.FastHappyColor().Clamped().RGB255()
	return []uint8{uint8(r), uint8(g), uint8(b)}
}
