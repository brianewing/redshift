package effects

import (
	"testing"
	"redshift/strip"
)

func animationSet() []Effect {
	return []Effect{
		&Clear{},
		&RainbowEffect{Size: 150, Speed: 1, Dynamic: true},
		//&Combine{
		//	Effects: []Effect{
		//		&Clear{},
		//		&RainbowEffect{Size: 100, Speed: 1, Dynamic: false},
		//		&Brightness{Brightness: 200},
		//	},
		//},
		&BlueEffect{},
		&LarsonEffect{Color: []uint8{0,0,0}},
	}
}

func stripAndEffects() (*strip.LEDStrip, []Effect) {
	strip := strip.New(60)
	effects := animationSet()

	for _, effect := range effects {
		effect.Render(strip)
	}

	return strip, effects
}

func BenchmarkMarshalJSON(b *testing.B) {
	_, effects := stripAndEffects()

	for i := 0; i < b.N; i++ {
		MarshalJson(effects)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	_, effects := stripAndEffects()
	effectsJson, _ := MarshalJson(effects)

	for i := 0; i < b.N; i++ {
		UnmarshalJson(effectsJson)
	}
}