package effects

import "github.com/brianewing/redshift/strip"

type Sepia struct {
	Factor float64
}

func (e *Sepia) Render(s *strip.LEDStrip) {
	for _, led := range s.Buffer {
		r, g, b := led[0], led[1], led[2]

		f := (255 - e.Factor) / 255

		led[0] = clamp(float64(r)*(0.393+0.607*f) + float64(g)*(0.769-0.769*f) + float64(b)*(0.189-0.189*f))
		led[1] = clamp(float64(r)*(0.349-0.349*f) + float64(g)*(0.686+0.314*f) + float64(b)*(0.168-0.168*f))
		led[2] = clamp(float64(r)*(0.272-0.272*f) + float64(g)*(0.534-0.534*f) + float64(b)*(0.131+0.869*f))
	}
}

func clamp(val float64) uint8 {
	if val > 255 {
		return 255
	} else {
		return uint8(val)
	}
}
