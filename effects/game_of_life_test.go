package effects

import (
	"testing"
)

func BenchmarkGameOfLifeHashing(b *testing.B) {
	g := newLife(50, 50, 500)

	b.Run("Step()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g.Step()
		}
	})

	b.Run("Step() -> append state to list", func(b *testing.B) {
		var states [][]bool

		for i := 0; i < b.N; i++ {
			g.Step()
			states = append(states, g.State())
		}
	})

	b.Run("Step() -> append -> compare", func(b *testing.B) {
		var states [][]bool

		for i := 0; i < b.N; i++ {
			g.Step()
			s := g.State()
			if len(states) > 200 {
				states = states[:200]
			}
			for _, s2 := range states {
				compareStates(s, s2)
			}
			states = append(states, s)
		}
	})

	// b.Run("Step() -> Hash()", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		g.Step()
	// 	}
	// })
}
