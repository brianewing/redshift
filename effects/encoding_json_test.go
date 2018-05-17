package effects

import (
	"testing"
)

func animationSet() EffectSet {
	return EffectSet{
		EffectEnvelope{Effect: &Clear{}},
		EffectEnvelope{Effect: &RainbowEffect{Size: 150, Speed: 1}},
		EffectEnvelope{Effect: &BlueEffect{}},
		EffectEnvelope{Effect: &LarsonEffect{Color: []uint8{0,0,0}}},
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	b.Run("Example", func(b *testing.B) {
		effects := animationSet()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MarshalJSON(effects)
		}
	})
	b.Run("Layer{Example}", func(b *testing.B) {
		effects := EffectSet{EffectEnvelope{Effect: &Layer{Effects: animationSet()}}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MarshalJSON(effects)
		}
	})
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("Example", func(b *testing.B) {
		effectsJson, _ := MarshalJSON(animationSet())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			UnmarshalJSON(effectsJson)
		}
	})
	b.Run("Layer{Example}", func(b *testing.B) {
		effectsJson, _ := MarshalJSON(EffectSet{EffectEnvelope{Effect: &Layer{Effects: animationSet()}}})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			UnmarshalJSON(effectsJson)
		}
	})
}
