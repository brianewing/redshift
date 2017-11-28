package effects

import (
	"github.com/brianewing/redshift/strip"
	"testing"
)

func BenchmarkBlending(b *testing.B) {
	s1 := strip.New(60)
	s1.Randomize()
	s2 := strip.New(s1.Size)
	s2.Randomize()

	b.ResetTimer()

	b.Run("RGB", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for i, c1 := range s1.Buffer {
				blendRgb(c1, s2.Buffer[i])
			}
		}
	})
	b.Run("HCL", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for i, c1 := range s1.Buffer {
				blendHcl(c1, s2.Buffer[i])
			}
		}
	})
}
