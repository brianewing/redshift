package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Effect interface {
	Render(strip *strip.LEDStrip)
}

type Initable interface {
	Init()
}
type Destroyable interface {
	Destroy()
}

type EffectEnvelope struct {
	Controls ControlSet
	Disabled bool
	Effect
}

func (e *EffectEnvelope) Init() {
	if initable, ok := e.Effect.(Initable); ok {
		initable.Init()
	}
	e.Controls.Init()
}

func (e *EffectEnvelope) Destroy() {
	if destroyable, ok := e.Effect.(Destroyable); ok {
		destroyable.Destroy()
	}
	e.Controls.Destroy()
}

func (e *EffectEnvelope) Render(strip *strip.LEDStrip) {
	if e.Disabled || e.Effect == nil {
		return
	}
	e.Controls.Apply(e.Effect)
	e.Effect.Render(strip)
}

type EffectSet []EffectEnvelope

func (s EffectSet) Init() {
	for _, envelope := range s {
		envelope.Init()
	}
}

func (s EffectSet) Destroy() {
	for _, envelope := range s {
		envelope.Destroy()
	}
}

func (s EffectSet) Render(strip *strip.LEDStrip) {
	for _, effect := range s {
		effect.Render(strip)
	}
}

/*
 * Construction
 */

func NewByName(name string) Effect {
	switch name {
	case "BlueEffect":
		return NewBlueEffect()
	case "Buffer", "Blend":
		return NewBlend()
	case "Brightness":
		return NewBrightness()
	case "Clear":
		return &Clear{}
	case "Channels":
		return NewChannels()
	case "External":
		return NewExternal()
	case "Fill":
		return NewFill()
	case "GameOfLife":
		return NewGameOfLife()
	case "Gamma":
		return NewGamma()
	case "Greyscale":
		return &Greyscale{}
	case "Layer":
		return NewLayer()
	case "Layout":
		return NewLayout()
	case "LarsonEffect":
		return NewLarsonEffect()
	case "Mirror":
		return NewMirror()
	case "MoodEffect":
		return NewMoodEffect()
	case "Rainbow", "RainbowEffect":
		return NewRainbow()
	case "RandomEffect":
		return &RandomEffect{}
	case "Script":
		return &Script{}
	case "Sepia":
		return NewSepia()
	case "Stripe":
		return NewStripe()
	case "Strobe":
		return NewStrobe()
	case "Slideshow":
		return NewSlideshow()
	case "Switch":
		return &Switch{}
	case "Trigger":
		return &Trigger{}
	case "Toggle":
		return NewToggle()
	case "Wheee", "Whoosh":
		return NewWhoosh()
	default:
		return &Null{}
	}
}

func Names() []string {
	return []string{
		"BlueEffect",
		"Brightness",
		"Buffer",
		"Clear",
		"Channels",
		"External",
		"Fill",
		"GameOfLife",
		"Gamma",
		"Greyscale",
		"Layer",
		"Layout",
		"LarsonEffect",
		"Mirror",
		"MoodEffect",
		"Null",
		"Rainbow",
		"RandomEffect",
		"Script",
		"Sepia",
		"Stripe",
		"Strobe",
		"Slideshow",
		"Switch",
		"Toggle",
		"Trigger",
		"Whoosh",
	}
}
