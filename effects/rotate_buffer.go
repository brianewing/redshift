package effects

import "redshift/strip"

type RotateBuffer struct {
	Buffer *[][]uint8
	Count int
	Reverse bool
}

func (e *RotateBuffer) Render(s *strip.LEDStrip) {
	*e.Buffer = rotate(*e.Buffer, e.Count, e.Reverse)
	(&Buffer{Buffer: *e.Buffer}).Render(s)
}

func rotate(buffer [][]uint8, n int, reverse bool) [][]uint8 {
	if reverse {
		head, tail := buffer[0:n], buffer[n:]
		return append(tail, head...)
	} else {
		head, tail := buffer[:len(buffer)-n], buffer[len(buffer)-n:]
		return append(tail, head...)
	}
}
