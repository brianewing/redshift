package server

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
)

// Redshift System Exclusive Messages

const OpcSystemId uint16 = 65535

type SystemExclusiveCmd uint8

const (
	CmdOpenStream = iota
	CmdCloseStream
	CmdSetStreamFps
	CmdSetStreamEffects
)

// Message format

type OpcMessage struct {
	Channel, Command uint8
	Length uint16
	Data []byte

	SystemExclusive SystemExclusive
}

type SystemExclusive struct {
	Command SystemExclusiveCmd
	Data []byte
}

func (se SystemExclusive) Bytes() []byte {
	bytes := append([]byte{byte(se.Command), 0, 0}, se.Data...)
	binary.BigEndian.PutUint16(bytes[1:3], uint16(len(se.Data)))
	return bytes
}

func (m OpcMessage) Bytes() []byte {
	bytes := []byte{m.Channel, m.Command, 0, 0}
	if m.Command == 255 {
		binary.BigEndian.PutUint16(bytes[2:4], OpcSystemId)
		return append(bytes, m.SystemExclusive.Bytes()...)
	} else {
		binary.BigEndian.PutUint16(bytes[2:4], m.Length)
		return append(bytes, m.Data...)
	}
}

func (m OpcMessage) WritePixels(buffer [][]uint8) {
	for i, val := range m.Data {
		if len(buffer) == i / 3 {
			break
		}
		buffer[i / 3][i % 3] = val
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

	if msg.Command == 255 { // system exclusive message
		if msg.Length == OpcSystemId {
			// sysex data format: [command, length high byte, length low byte, data..]
			// i.e. opc without leading channel byte

			sysHeader := make([]byte, 3)

			if bytesRead, err := r.Read(sysHeader); err != nil {
				return msg, err
			} else if bytesRead != 3 {
				return msg, errors.New("bad header")
			}

			msg.SystemExclusive.Command = SystemExclusiveCmd(sysHeader[0])
			log.Println("sysex header", sysHeader)

			sysDataLength := binary.BigEndian.Uint16(sysHeader[1:3])
			msg.SystemExclusive.Data = make([]byte, sysDataLength)

			if bytesRead, err := r.Read(msg.SystemExclusive.Data); err != nil {
				return msg, err
			} else if bytesRead != int(sysDataLength) {
				return msg, errors.New("system exclusive data length mismatch")
			} else {
				log.Println("got", bytesRead, "bytes of sys data")
				return msg, nil
			}
		} else {
			return msg, errors.New("system exclusive system id mismatch")
		}
	}

	msg.Data = make([]byte, msg.Length)
	bytesRead, err = r.Read(msg.Data)

	if err != nil {
		return msg, err
	} else if bytesRead != int(msg.Length) {
		return msg, errors.New("data length mismatch")
	} else {
		return msg, nil
	}
}
