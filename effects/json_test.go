package effects

import (
	"testing"
	"redshift/strip"
)

func BenchmarkMarshal(b *testing.B) {
	strip := strip.New(60)

	effects := []Effect{
		&Clear{},
		&Brightness{Brightness: 200},
		&RainbowEffect{Size: 150, Speed: 0},
	}

	for _, effect := range effects {
		effect.Render(strip)
	}

	for i := 0; i < b.N; i++ {
		MarshalJson(effects)
	}
}