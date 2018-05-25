package effects

import (
	"github.com/brianewing/redshift/strip"
)

type Greyscale struct{}

func (e *Greyscale) Render(s *strip.LEDStrip) {
	for _, led := range s.Buffer {
		copy(led, convertToGreyscale(led))
	}
}

func convertToGreyscale(led []uint8) []uint8 {
	r := float64(led[0]) / 255
	g := float64(led[1]) / 255
	b := float64(led[2]) / 255
	
	y := uint8((0.299 * r + 0.587 * g + 0.114 * b) * 255)
	return []uint8{y, y, y}
}
