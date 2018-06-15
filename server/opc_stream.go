package server

import (
	"encoding/json"
	"github.com/brianewing/redshift/animator"
	"time"
)

type opcStream struct {
	channel  uint8
	animator *animator.Animator

	stop             chan struct{}
	fpsChange        chan uint8
	effectsFpsChange chan uint8

	virtual bool
}

func NewOpcStream(channel uint8) *opcStream {
	return &opcStream{
		channel:          channel,
		stop:             make(chan struct{}),
		fpsChange:        make(chan uint8),
		effectsFpsChange: make(chan uint8),
	}
}

func (s *opcStream) Close() {
	s.stop <- struct{}{}
}

func (s *opcStream) SetFps(fps uint8) {
	s.fpsChange <- fps
}

func (s *opcStream) SetEffectsFps(fps uint8) {
	s.effectsFpsChange <- fps
}

func (s *opcStream) Run(w OpcWriter) {
	var pixelTicker *time.Ticker
	var pixelChan <-chan time.Time

	var effectsTicker *time.Ticker
	var effectsChan <-chan time.Time

	if s.virtual {
		go s.animator.Run(16 * time.Millisecond) // 60 fps
	}

	for {
		select {
		case <-pixelChan:
			s.WritePixels(w)

		case <-effectsChan:
			s.WriteEffects(w)

		case newFps := <-s.fpsChange:
			if pixelTicker != nil {
				pixelTicker.Stop()
			}
			if newFps > 0 {
				pixelTicker = time.NewTicker(time.Second / time.Duration(newFps))
				pixelChan = pixelTicker.C
			} else {
				s.WritePixels(w)
			}

		case newEffectsFps := <-s.effectsFpsChange:
			if effectsTicker != nil {
				effectsTicker.Stop()
			}
			if newEffectsFps > 0 {
				effectsTicker = time.NewTicker(time.Second / time.Duration(newEffectsFps))
				effectsChan = effectsTicker.C
			} else {
				s.WriteEffects(w)
			}

		case <-s.stop:
			if pixelTicker != nil {
				pixelTicker.Stop()
			}
			if effectsTicker != nil {
				effectsTicker.Stop()
			}
			if s.virtual {
				s.animator.Finish()
			}
			return
		}
	}
}

func (s *opcStream) WritePixels(w OpcWriter) error {
	s.animator.Strip.Lock()
	pixels := s.animator.Strip.MarshalBytes()
	s.animator.Strip.Unlock()
	msg := OpcMessage{
		Channel: s.channel,
		Command: 0, // write pixels
		Length:  uint16(len(pixels)),
		Data:    pixels,
	}
	return w.WriteOpc(msg)
}

func (s *opcStream) WriteEffects(w OpcWriter) error {
	if effectsJson, err := json.Marshal(s.animator.Effects); err != nil {
		return err
	} else {
		msg := OpcMessage{
			Channel: s.channel,
			Command: 255, // system exclusive
			SystemExclusive: SystemExclusive{
				Command: CmdSetEffectsJson,
				Data:    []byte(effectsJson),
			},
		}
		return w.WriteOpc(msg)
	}
}
