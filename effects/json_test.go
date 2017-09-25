package effects

import (
	"testing"
)

func animationSet() []Effect {
	return []Effect{
		&Clear{},
		&RainbowEffect{Size: 150, Speed: 1, Dynamic: true},
		&BlueEffect{},
		&LarsonEffect{Color: []uint8{0,0,0}},
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	b.Run("Example", func(b *testing.B) {
		effects := animationSet()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MarshalJson(effects)
		}
	})
	b.Run("Combine{Example}", func(b *testing.B) {
		effects := []Effect{&Combine{Effects: animationSet()}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MarshalJson(effects)
		}
	})
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("Example", func(b *testing.B) {
		effectsJson, _ := MarshalJson(animationSet())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			UnmarshalJson(effectsJson)
		}
	})
	b.Run("Combine{Example}", func(b *testing.B) {
		effectsJson, _ := MarshalJson([]Effect{&Combine{Effects: animationSet()}})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			UnmarshalJson(effectsJson)
		}
	})
}
