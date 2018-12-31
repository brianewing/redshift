package opc

import (
	"encoding/json"
	"github.com/brianewing/redshift/animator"
	"time"
)

type stream struct {
	channel  uint8
	animator *animator.Animator

	stop             chan struct{}
	fpsChange        chan uint8
	effectsFpsChange chan uint8

	virtual bool
}

func newStream(channel uint8) *stream {
	return &stream{
		channel:          channel,
		stop:             make(chan struct{}),
		fpsChange:        make(chan uint8),
		effectsFpsChange: make(chan uint8),
	}
}

func (s *stream) Close() {
	s.stop <- struct{}{}
}

func (s *stream) SetFps(fps uint8) {
	s.fpsChange <- fps
}

func (s *stream) SetEffectsFps(fps uint8) {
	s.effectsFpsChange <- fps
}

func (s *stream) Run(w Writer) {
	var pixelTicker *time.Ticker
	var pixelChan <-chan time.Time

	var effectsTicker *time.Ticker
	var effectsChan <-chan time.Time

	if s.virtual {
		go s.animator.Run(time.Second / 60) // 60 fps
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

func (s *stream) WritePixels(w Writer) error {
	s.animator.Strip.Lock()
	pixels := s.animator.Strip.MarshalBytes()
	s.animator.Strip.Unlock()
	msg := Message{
		Channel: s.channel,
		Command: 0, // write pixels
		Length:  uint16(len(pixels)),
		Data:    pixels,
	}
	return w.WriteOpc(msg)
}

func (s *stream) WriteEffects(w Writer) error {
	if effectsJson, err := json.Marshal(s.animator.Effects); err != nil {
		return err
	} else {
		msg := Message{
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
