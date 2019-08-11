package opc

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/osc"
	"github.com/brianewing/redshift/repl"
	"github.com/brianewing/redshift/strip"
)

var REDSHIFT_VERSION = "0.1.0"

type Writer interface {
	WriteOpc(msg Message) error
}

type Session struct {
	Animator *animator.Animator
	Client   Writer
	ClientInfo

	streams map[uint8]*stream   // opened, indexed by channel
	repls   map[uint8]*replPipe // indexed by channel
	sync.Mutex
}

func NewSession(animator *animator.Animator, client Writer) *Session {
	return &Session{
		Animator: animator,
		Client:   client,

		streams: make(map[uint8]*stream),
		repls:   make(map[uint8]*replPipe),
	}
}

func (s *Session) Receive(msg Message) error {
	s.Lock()

	switch msg.Command {
	case 0:
		if s.Animator != nil {
			s.Animator.Strip.Lock()
			msg.WritePixels(s.Animator.Strip.Buffer)
			s.Animator.Strip.Unlock()
		}

	case 255:
		switch msg.SystemExclusive.Command {
		// confirms successful connection by sending server info
		case CmdWelcome:
			s.receiveClientInfo(msg.SystemExclusive.Data)
			s.sendWelcome()

		// identifies which osc addresses have received msgs so far
		case CmdOscSummary:
			s.sendOscSummary(msg.Channel)

		// clears the list of addresses that have received msgs so far
		case CmdClearOscSummary:
			osc.ClearSummary()
			s.sendOscSummary(msg.Channel)

		case CmdOpenStream:
			channel := msg.Channel
			description := string(msg.SystemExclusive.Data)

			if stream, err := s.openStream(channel, description); err == nil {
				s.streams[msg.Channel] = stream
			} else {
				s.sendError(msg.Channel, CmdOpenStream, err)
				return err
			}

		case CmdCloseStream:
			stream := s.streams[msg.Channel]
			stream.Close()
			s.streams[msg.Channel] = nil

		case CmdSetStreamFps:
			stream := s.streams[msg.Channel]
			fps := msg.SystemExclusive.Data[0]
			stream.SetFps(fps)

		case CmdSetEffectsStreamFps:
			stream := s.streams[msg.Channel]
			fps := msg.SystemExclusive.Data[0]
			stream.SetEffectsFps(fps)

		case CmdSetEffectsJson:
			if newEffects, err := effects.UnmarshalJSON(msg.SystemExclusive.Data); err != nil {
				s.sendError(msg.Channel, CmdSetEffectsJson, err)
			} else {
				stream := s.streams[msg.Channel]
				stream.animator.SetEffects(newEffects)
				stream.WriteEffects(s.Client)
			}

		case CmdSetEffectsYaml:
			if newEffects, err := effects.UnmarshalYAML(msg.SystemExclusive.Data); err != nil {
				s.sendError(msg.Channel, CmdSetEffectsYaml, err)
			} else {
				stream := s.streams[msg.Channel]
				stream.animator.SetEffects(newEffects)
			}

		case CmdAppendEffectsJson:
			if newEffects, err := effects.UnmarshalJSON(msg.SystemExclusive.Data); err != nil {
				s.sendError(msg.Channel, CmdAppendEffectsJson, err)
			} else {
				newEffects.Init()
				stream := s.streams[msg.Channel]
				stream.animator.Effects = append(stream.animator.Effects, newEffects...)
			}

		case CmdAppendEffectsYaml:
			if newEffects, err := effects.UnmarshalYAML(msg.SystemExclusive.Data); err != nil {
				s.sendError(msg.Channel, CmdAppendEffectsYaml, err)
			} else {
				newEffects.Init()
				stream := s.streams[msg.Channel]
				stream.animator.Effects = append(stream.animator.Effects, newEffects...)
			}

		case CmdRepl:
			if s.repls[msg.Channel] == nil {
				pipe := newReplPipe(msg.Channel, s.Client)
				s.repls[msg.Channel] = &pipe
				go repl.Run(s.Animator, s.repls[msg.Channel], "> ") // will exit when repl pipe is closed
			}

			s.repls[msg.Channel].inputWriter.Write(append(msg.SystemExclusive.Data, '\n'))

		case CmdPing:
			s.sendPong(msg.Channel, msg.SystemExclusive.Data)

		case CmdClose:
			if c, ok := s.Client.(io.Closer); ok {
				c.Close()
			}

		default:
			println("unrecognised opc system cmd", strconv.Itoa(int(msg.SystemExclusive.Command)))
		}

	default:
		return errors.New("command not recognised")
	}

	s.Unlock()
	return nil
}

var startTime = time.Now()

func (s *Session) sendWelcome() error {
	welcomeJson, _ := json.Marshal(map[string]interface{}{
		"version": REDSHIFT_VERSION,
		"started": startTime,
		"config": map[string]interface{}{
			"serverName": "Living Room Ceiling Strip",
		},
		"availableEffects": effects.Names(),
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

func (s *Session) sendPong(channel uint8, data []byte) error {
	return s.Client.WriteOpc(Message{
		Channel: channel,
		Command: 255,
		SystemExclusive: SystemExclusive{
			Command: CmdPong,
			Data:    data,
		},
	})
}

func (s *Session) openStream(channel uint8, description string) (*stream, error) {
	if s.streams[channel] != nil {
		return nil, errors.New("channel has an open stream, please close before opening another")
	}

	stream := makeStream(channel)
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

func (s *Session) receiveClientInfo(data []byte) {
	if len(data) > 0 {
		json.Unmarshal(data, &s.ClientInfo)
		log.Println("client info", s.ClientInfo)
	}
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

func (s *Session) sendError(channel uint8, cmd SystemExclusiveCmd, err error) {
	if err != nil {
		log.Println(s, "Error", "|", err)
		s.Client.WriteOpc(Message{
			Channel: channel,
			Command: 255,
			SystemExclusive: SystemExclusive{
				Command: CmdErrorOccurred,
				Data:    append([]byte{byte(cmd)}, err.Error()...),
			},
		})
	}
}

func (s *Session) Close() {
	for _, stream := range s.streams {
		stream.Close()
	}
	for _, replSession := range s.repls {
		replSession.Close()
	}

}

func (s *Session) String() string {
	return "Session{" + s.ClientInfo.DeviceName + "}"
}
