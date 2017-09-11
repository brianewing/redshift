package effects

import "redshift/strip"

type Effect interface {
	Render(strip *strip.LEDStrip)
}

