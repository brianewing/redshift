package server

import (
	"github.com/brianewing/redshift/animator"
	"time"
)

type opcStream struct {
	channel  uint8
	animator *animator.Animator

	fps        uint8
	fpsChanged chan struct{}

	virtual bool // if true the animation will be stopped when the stream is closed
}

func NewOpcStream(channel uint8) *opcStream {
	return &opcStream{fpsChanged: make(chan struct{})}
}

func (s *opcStream) SetFps(fps uint8) {
	s.fps = fps
	s.fpsChanged <- struct{}{}
}

func (s *opcStream) Run(w OpcWriter) {
	var ticker *time.Ticker
	var C <-chan time.Time

	println("run opcstream")

	for {
		select {
		case <-s.fpsChanged:
			if ticker != nil {
				ticker.Stop()
			}
			if s.fps > 0 {
				ticker = time.NewTicker(time.Second / time.Duration(s.fps))
				C = ticker.C
			} else {
				C = nil
			}
		case <-C:
			s.WriteFrame(w)
		}
	}
}

func (s *opcStream) WriteFrame(w OpcWriter) error {
	s.animator.Strip.Lock()
	pixels := s.animator.Strip.SerializeBytes()
	msg := OpcMessage{
		Channel: s.channel,
		Command: 0, // write pixels
		Length: uint16(len(pixels)),
		Data: pixels,
	}
	s.animator.Strip.Unlock()
	return w.WriteOpc(msg)
}
