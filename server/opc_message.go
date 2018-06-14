package server

import (
	"github.com/brianewing/redshift/strip"
	"encoding/binary"
	"errors"
	"io"
)

type OpcMessage struct {
	Channel, Command uint8
	Length           uint16
	Data             []byte

	SystemExclusive SystemExclusive
}

func (m OpcMessage) Bytes() []byte {
	bytes := []byte{m.Channel, m.Command, 0, 0}
	if m.Command == 255 {
		sysex := m.SystemExclusive.Bytes()
		m.Length = uint16(len(sysex))
		m.Data = sysex
	}
	binary.BigEndian.PutUint16(bytes[2:4], m.Length)
	return append(bytes, m.Data...)
}

func (m OpcMessage) WritePixels(buffer strip.Buffer) {
	for i, val := range m.Data {
		if len(buffer) == i/3 {
			break
		}
		buffer[i/3][i%3] = val
	}
}

func ReadOpcMessage(r io.Reader) (OpcMessage, error) {
	// message format : [channel, command, length high byte, length low byte, data...]
	msg := OpcMessage{}

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
	bytesRead, err = r.Read(msg.Data)

	if msg.Command == 255 {
		if len(msg.Data) < 2 {
			return msg, errors.New("missing sysex data")
		}

		msg.SystemExclusive.SystemId = binary.BigEndian.Uint16(msg.Data[0:2])

		if msg.SystemExclusive.SystemId == OpcSystemId {
			if len(msg.Data) < 3 {
				return msg, errors.New("missing sysex command byte")
			}
			msg.SystemExclusive.Command = SystemExclusiveCmd(msg.Data[2])
			msg.SystemExclusive.Data = msg.Data[3:]
		} else {
			msg.SystemExclusive.Data = msg.Data[2:]
		}
	}

	if err != nil {
		return msg, err
	} else if bytesRead != int(msg.Length) {
		return msg, errors.New("data length mismatch")
	} else {
		return msg, nil
	}
}

// System exclusive messages

const OpcSystemId uint16 = 65535

type SystemExclusiveCmd uint8

const (
	CmdWelcome = iota

	CmdOpenStream
	CmdCloseStream
	CmdSetStreamFps

	CmdSetEffectsJson
	CmdSetEffectsStreamFps
	CmdSetEffectsYaml

	CmdAppendEffectsJson
	CmdAppendEffectsYaml
)

type SystemExclusive struct {
	SystemId uint16
	Command  SystemExclusiveCmd
	Data     []byte
}

func (se SystemExclusive) Bytes() []byte {
	bytes := append([]byte{0, 0, byte(se.Command)}, se.Data...)
	binary.BigEndian.PutUint16(bytes[0:2], OpcSystemId)
	return bytes
}
