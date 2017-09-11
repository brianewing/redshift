package effects

import "redshift/strip"

type Buffer struct {
	Buffer [][]int
}

func (e *Buffer) Render(strip *strip.LEDStrip) {
	for i := range e.Buffer {
		copy(strip.Buffer[i], e.Buffer[i])
	}
}
