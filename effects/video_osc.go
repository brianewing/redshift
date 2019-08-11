package effects

import (
	"github.com/brianewing/redshift/strip"
	"github.com/brianewing/redshift/osc"

	"strconv"
	"strings"
)

// VideoOSC is an Effect which copies pixel values from OSC messages in the form used by the VideoOSC Android app, e.g.
//	/vosc/red1 [255] /vosc/green1 [255] /vosc/blue1 [255]
//	/vosc/red2 [255] /vosc/green2 [255] /vosc/blue2 [255]
// onto LED strip buffer, where /vosc/red1 specifies the "0th" LED color's red component,
// and so on.
type VideoOSC struct {
	Prefix string
	Blend *Blend

	stop chan struct{}
}

func NewVideoOSC() *VideoOSC {
	return &VideoOSC{
		Prefix: "/vosc/",
		Blend: NewBlend(),
	}
}

func (e *VideoOSC) Init() {
	e.stop = make(chan struct{})

	go func() {
		msgs, done := osc.StreamMessages()

		for {
			select {
			case msg := <-msgs:
				e.receiveOSCMessage(msg)
			case <-e.stop:
				done <- struct{}{}
				return
			}
		}
	}()
}

func (e *VideoOSC) Destroy() {
	e.stop <- struct{}{}
}

func (e *VideoOSC) Render(s *strip.LEDStrip) {
	if len(e.Blend.Buffer) != int(s.Size) {
		e.Blend.Buffer = strip.NewBuffer(s.Size)
	}

	e.Blend.Render(s)
}

func (e *VideoOSC) receiveOSCMessage(msg osc.OscMessage) {
	if desc := strings.TrimPrefix(msg.Address, e.Prefix); len(msg.Address) > len(desc) {
		if len(msg.Arguments) > 0 {
			value, _ := msg.Arguments[0].(float64)

			var ledNum int
			var component uint8

			if numStr := strings.TrimPrefix(desc, "red"); len(desc) > len(numStr) {
				ledNum, _ = strconv.Atoi(numStr)
			} else if numStr := strings.TrimPrefix(desc, "green"); len(desc) > len(numStr) {
				ledNum, _ = strconv.Atoi(numStr)
				component = 1
			} else if numStr := strings.TrimPrefix(desc, "blue"); len(desc) > len(numStr) {
				ledNum, _ = strconv.Atoi(numStr)
				component = 2
			}

			ledNum-- // VideoOSC indexes pixels from 1
			e.updatePixelComponent(ledNum, component, uint8(value))
		}
	}
}

func (e *VideoOSC) updatePixelComponent(pixel int, component, value uint8) {
	if pixel >= len(e.Blend.Buffer) {
		return
	}

	e.Blend.Buffer[pixel][component] = value
}