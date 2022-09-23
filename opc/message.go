package opc

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/brianewing/redshift/strip"
)

type Message struct {
	Channel, Command uint8
	Length           uint16
	Data             []byte

	SystemExclusive SystemExclusive
}

func (m Message) Bytes() []byte {
	bytes := []byte{m.Channel, m.Command, 0, 0}
	if m.Command == 255 {
		sysex := m.SystemExclusive.Bytes()
		m.Length = uint16(len(sysex))
		m.Data = sysex
	} else {
		m.Length = uint16(len(m.Data))
	}
	binary.BigEndian.PutUint16(bytes[2:4], m.Length)
	return append(bytes, m.Data...)
}

func (m Message) WritePixels(buffer strip.Buffer) {
	for i, val := range m.Data {
		if len(buffer) == i/3 {
			break
		}
		buffer[i/3][i%3] = val
	}
}

func ReadMessage(r io.Reader) (Message, error) {
	// message format : [channel, command, length high byte, length low byte, data...]
	msg := Message{}

	header := make([]byte, 4)
	bytesRead, err := r.Read(header)

	if err != nil {
		return msg, err
	} else if bytesRead != 4 {
		return msg, errors.New("bad header")
	}

	msg.Channel = header[0]
	msg.Command = header[1]
	msg.Length = binary.BigEndian.Uint16(header[2:4])

	msg.Data = make([]byte, msg.Length)

	if bytesRead, err = r.Read(msg.Data); err != nil {
		return msg, err
	} else if bytesRead != int(msg.Length) {
		return msg, errors.New("data length mismatch")
	}

	if msg.Command == 255 {
		if len(msg.Data) < 2 {
			return msg, errors.New("missing sysex system id")
		} else if len(msg.Data) < 3 {
			return msg, errors.New("missing sysex command")
		}

		msg.SystemExclusive.SystemID = binary.BigEndian.Uint16(msg.Data[0:2])

		if msg.SystemExclusive.SystemID == OpcSystemID {
			msg.SystemExclusive.Command = SystemExclusiveCmd(msg.Data[2])
			msg.SystemExclusive.Data = msg.Data[3:]
		} else {
			msg.SystemExclusive.Data = msg.Data[2:]
		}
	}

	return msg, nil
}

// System exclusive messages

const OpcSystemID uint16 = 65535

type SystemExclusiveCmd uint8

const (
	CmdWelcome SystemExclusiveCmd = iota

	CmdOpenStream
	CmdCloseStream
	CmdSetStreamFps

	CmdSetEffectsJson
	CmdSetEffectsStreamFps
	CmdSetEffectsYaml

	CmdAppendEffectsJson
	CmdAppendEffectsYaml

	CmdOscSummary
	CmdClearOscSummary

	CmdErrorOccurred

	CmdRepl

	CmdPing
	CmdPong

	CmdClose

	CmdTickAnimation
)

type SystemExclusive struct {
	SystemID uint16
	Command  SystemExclusiveCmd
	Data     []byte
}

func (se SystemExclusive) Bytes() []byte {
	bytes := append([]byte{0, 0, byte(se.Command)}, se.Data...)
	if se.SystemID == 0 {
		binary.BigEndian.PutUint16(bytes[0:2], OpcSystemID)
	} else {
		binary.BigEndian.PutUint16(bytes[0:2], se.SystemID)
	}
	return bytes
}
