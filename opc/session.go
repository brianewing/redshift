package opc

import (
	"encoding/json"
	"errors"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/osc"
	"github.com/brianewing/redshift/strip"
	"strconv"
	"strings"
	"time"
)

var REDSHIFT_VERSION = "0.1.0"

type Writer interface {
	WriteOpc(msg Message) error
}

type Session struct {
	Animator *animator.Animator
	Client   Writer

	streams []*opcStream
}

func (s *Session) Receive(msg Message) error {
	switch msg.Command {
	case 0:
		if s.Animator != nil {
			s.Animator.Strip.Lock()
			msg.WritePixels(s.Animator.Strip.Buffer)
			s.Animator.Strip.Unlock()
		}
	case 255:
		switch msg.SystemExclusive.Command {
		case CmdWelcome:
			s.sendWelcome() // confirms successful connection by sending server info
		case CmdOscSummary:
			s.sendOscSummary(msg.Channel) // identifies which osc addresses are receiving msgs
		case CmdClearOscSummary:
			osc.ClearSummary()
			s.sendOscSummary(msg.Channel)
		case CmdOpenStream:
			channel := msg.Channel
			description := string(msg.SystemExclusive.Data)

			if stream, err := s.openStream(channel, description); err == nil {
				s.streams = append(s.streams, stream)
			} else {
				return err
			}
		case CmdCloseStream:
			stream := s.streams[msg.Channel]
			stream.Close()
			s.streams = append(s.streams[:msg.Channel], s.streams[msg.Channel+1:]...)
		case CmdSetStreamFps:
			stream := s.streams[msg.Channel]
			fps := msg.SystemExclusive.Data[0]
			stream.SetFps(fps)
		case CmdSetEffectsStreamFps:
			stream := s.streams[msg.Channel]
			fps := msg.SystemExclusive.Data[0]
			stream.SetEffectsFps(fps)
		case CmdSetEffectsJson:
			newEffects, _ := effects.UnmarshalJSON(msg.SystemExclusive.Data)
			stream := s.streams[msg.Channel]
			stream.animator.SetEffects(newEffects)
		case CmdSetEffectsYaml:
			newEffects, _ := effects.UnmarshalYAML(msg.SystemExclusive.Data)
			stream := s.streams[msg.Channel]
			stream.animator.SetEffects(newEffects)
		case CmdAppendEffectsJson:
			newEffects, _ := effects.UnmarshalJSON(msg.SystemExclusive.Data)
			newEffects.Init()
			stream := s.streams[msg.Channel]
			stream.animator.Effects = append(stream.animator.Effects, newEffects...)
		case CmdAppendEffectsYaml:
			newEffects, _ := effects.UnmarshalYAML(msg.SystemExclusive.Data)
			newEffects.Init()
			stream := s.streams[msg.Channel]
			stream.animator.Effects = append(stream.animator.Effects, newEffects...)
		default:
			println("dont know how to handle system cmd", strconv.Itoa(int(msg.SystemExclusive.Command)))
		}
	default:
		return errors.New("command not recognised")
	}
	return nil
}

var startTime = time.Now()

func (s *Session) sendWelcome() error {
	welcomeJson, _ := json.Marshal(map[string]interface{}{
		"version": REDSHIFT_VERSION,
		"started": startTime,
	})
	msg := Message{
		Command: 255,
		SystemExclusive: SystemExclusive{
			Command: CmdWelcome,
			Data:    welcomeJson,
		},
	}
	return s.Client.WriteOpc(msg)
}

func (s *Session) openStream(channel uint8, description string) (*opcStream, error) {
	stream := NewOpcStream(channel)
	desc := strings.Fields(description)

	switch desc[0] {
	case "strip", "":
		stream.animator = s.Animator
	case "virtual":
		stream.virtual = true
		stream.animator = &animator.Animator{}

		numLeds := 30
		if len(desc) >= 2 {
			if v, _ := strconv.Atoi(desc[1]); v > 0 {
				numLeds = v
			}
		}

		stream.animator.Strip = strip.New(numLeds)
	}

	go stream.Run(s.Client)
	return stream, nil
}

func (s *Session) sendOscSummary(channel uint8) {
	oscMsgs := osc.Summary()
	jsonBytes, _ := json.Marshal(oscMsgs)

	s.Client.WriteOpc(Message{
		Channel: channel,
		Command: 255,
		SystemExclusive: SystemExclusive{
			Command: CmdOscSummary,
			Data:    jsonBytes,
		},
	})
}

func (s *Session) Close() {
	for _, stream := range s.streams {
		stream.Close()
	}
}
