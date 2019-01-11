package effects

import "github.com/brianewing/redshift/strip"

type Sepia struct{}

func (e *Sepia) Render(s *strip.LEDStrip) {
	for _, led := range s.Buffer {
		r, g, b := led[0], led[1], led[2]

		led[0] = clamp(float64(r)*0.393 + float64(g)*0.769 + float64(b)*0.189)
		led[1] = clamp(float64(r)*0.349 + float64(g)*0.686 + float64(b)*0.168)
		led[2] = clamp(float64(r)*0.272 + float64(g)*0.534 + float64(b)*0.131)
	}
}

func clamp(val float64) uint8 {
	if val > 255 {
		return 255
	} else {
		return uint8(val)
	}
}
