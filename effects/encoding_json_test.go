package effects

import (
	"testing"
)

func animationSet() []Effect {
	return []Effect{
		&Clear{},
		&RainbowEffect{Size: 150, Speed: 1},
		&BlueEffect{},
		&LarsonEffect{Color: []uint8{0,0,0}},
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
		effects := []Effect{&Layer{Effects: animationSet()}}
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
		effectsJson, _ := MarshalJSON([]Effect{&Layer{Effects: animationSet()}})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			UnmarshalJSON(effectsJson)
		}
	})
}
