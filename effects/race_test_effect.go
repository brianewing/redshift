package effects

import (
	"redshift/strip"
)

type RaceTestEffect struct {}

// Continuously wipes and sets every pixel to red...
// A race condition exists if a client ever receives black

func (e *RaceTestEffect) Render(s *strip.LEDStrip) {
	s.Clear()
	for i := range s.Buffer {
		s.Buffer[i] = []int{255, 0, 0}
	}
}
