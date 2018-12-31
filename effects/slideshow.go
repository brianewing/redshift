package effects

import "github.com/brianewing/redshift/strip"

type Slideshow struct {
	Speed float64
	ActiveEffect EffectSet
	switcher Switch
}

func NewSlideshow() *Slideshow {
	return &Slideshow{
		Speed: 0.03,
	}
}

func (e *Slideshow) Init() {
	for _, name := range Names() {
		if name != "Slideshow" {
			e.switcher.Effects = append(e.switcher.Effects, EffectEnvelope{Effect: NewByName(name)})
		}
	}
	e.switcher.Init()
}

func (e *Slideshow) Destroy() {
	e.switcher.Destroy()
}

func (e *Slideshow) Render(s *strip.LEDStrip) {
	e.switcher.Selection = round(SmoothOscillateBetween(0, float64(len(e.switcher.Effects)-1), e.Speed))
	e.ActiveEffect = EffectSet{e.switcher.Effects[e.switcher.Selection]}
	e.switcher.Render(s)
}
