package effects

import "github.com/brianewing/redshift/strip"
import "github.com/longears/pixelslinger/colorutils"

type Gamma struct {
	Value uint8
}

func NewGamma() *Gamma {
	return &Gamma{Value: 127}
}

func (e *Gamma) Render(strip *strip.LEDStrip) {
	for _, led := range strip.Buffer {
		r, g, b := colorutils.GammaRgb(
			float64(led[0])/255,
			float64(led[1])/255,
			float64(led[2])/255,

			1/(float64(e.Value)/255*2),
		)

		led[0], led[1], led[2] = uint8(r*255), uint8(g*255), uint8(b*255)
	}
}
