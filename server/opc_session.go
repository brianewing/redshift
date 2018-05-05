package server

import (
	"encoding/json"
	"errors"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"log"
	"strconv"
	"strings"
)

type OpcWriter interface {
	WriteOpc(msg OpcMessage) error
}

type OpcSession struct {
	client OpcWriter

	animator *animator.Animator
	streams  []*opcStream
}

func (s *OpcSession) Receive(msg OpcMessage) error {
	switch msg.Command {
	case 0:
		log.Println("got incoming pixels", msg)
	case 255:
		log.Println("got sysex", msg)
		switch msg.SystemExclusive.Command {
		case CmdOpenStream:
			channel := msg.Channel
			description := string(msg.SystemExclusive.Data)

			log.Println("open stream", channel, description)

			if stream, err := s.openStream(channel, description); err == nil {
				s.streams = append(s.streams, stream)
			} else {
				return err
			}
		case CmdSetStreamFps:
			stream := s.streams[msg.Channel]
			fps := msg.SystemExclusive.Data[0]
			stream.SetFps(fps)
		case CmdSetEffectsStreamFps:
			stream := s.streams[msg.Channel]
			fps := msg.SystemExclusive.Data[0]
			stream.SetEffectsFps(fps)
		case CmdSetEffectsJson:
			var newEffects effects.EffectSet
			json.Unmarshal(msg.SystemExclusive.Data, &newEffects)
			stream := s.streams[msg.Channel]
			stream.animator.SetEffects(newEffects)
		default:
			println("dont know how to handle system cmd", strconv.Itoa(int(msg.SystemExclusive.Command)))
		}
	default:
		return errors.New("command not recognised")
	}
	return nil
}

func (s *OpcSession) openStream(channel uint8, description string) (*opcStream, error) {
	stream := NewOpcStream(channel)

	switch description {
	case "strip":
		stream.animator = s.animator
	case "virtual":
		stream.virtual = true
		stream.animator = &animator.Animator{} // todo: run + stop on close
	}

	go stream.Run(s.client)
	return stream, nil
}

func (s *OpcSession) Close() {}
