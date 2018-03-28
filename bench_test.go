package main

import (
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"testing"
)

func BenchmarkExampleAnimation(b *testing.B) {
	ledStrip := strip.New(60)

	randomStrip := strip.New(ledStrip.Size)
	randomStrip.Randomize()

	animator := &animator.Animator{
		Strip: ledStrip,
		Effects: effects.EffectSet{
			effects.EffectEnvelope{Effect: &effects.Clear{}},
			effects.EffectEnvelope{Effect: &effects.Blend{Buffer: randomStrip.Buffer}},
			effects.EffectEnvelope{Effect: &effects.RainbowEffect{Size: 150, Speed: 1}},
			effects.EffectEnvelope{Effect: &effects.Layer{
				Effects: effects.EffectSet{
					effects.EffectEnvelope{Effect: &effects.Clear{}},
					effects.EffectEnvelope{Effect: &effects.RainbowEffect{Size: 100, Speed: 1}},
					effects.EffectEnvelope{Effect: &effects.Brightness{Brightness: 200}},
				},
			}},
			effects.EffectEnvelope{Effect: &effects.BlueEffect{}},
			effects.EffectEnvelope{Effect: &effects.LarsonEffect{Color: []uint8{0, 0, 0}}},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		animator.Render()
	}
}
