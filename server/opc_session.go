package server

import (
	"encoding/json"
	"errors"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"log"
	"strconv"
	"strings"
	"time"
)

var REDSHIFT_VERSION = "0.1.0"

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
		switch msg.SystemExclusive.Command {
		case CmdWelcome:
			s.sendWelcome() // confirms successful connection by sending server info
		case CmdOpenStream:
			channel := msg.Channel
			description := string(msg.SystemExclusive.Data)

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

var startTime = time.Now()

func (s *OpcSession) sendWelcome() error {
	welcomeJson, _ := json.Marshal(map[string]interface{}{
		"version": REDSHIFT_VERSION,
		"started": startTime,
		"uptime":  time.Now().Sub(startTime).Seconds(),
	})
	msg := OpcMessage{
		Command: 255,
		SystemExclusive: SystemExclusive{
			Command: CmdWelcome,
			Data:    welcomeJson,
		},
	}
	return s.client.WriteOpc(msg)
}

func (s *OpcSession) openStream(channel uint8, description string) (*opcStream, error) {
	stream := NewOpcStream(channel)
	desc := strings.Fields(description)

	switch desc[0] {
	case "strip":
		stream.animator = s.animator
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

	go stream.Run(s.client)
	return stream, nil
}

func (s *OpcSession) Close() {
	for _, stream := range s.streams {
		stream.Close()
	}
}