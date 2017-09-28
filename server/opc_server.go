package server

import (
	"net"
	"log"
	"errors"
	"encoding/binary"
	"redshift/strip"
)

type OpcServer struct {
	Messages chan *OpcMessage
}

func RunOpcServer(addr string, strip *strip.LEDStrip) error {
	s := &OpcServer{Messages: make(chan *OpcMessage)}

	go func() {
		for {
			if msg := <-s.Messages; msg.Command == 0 {
				msg.WritePixels(strip.Buffer)
			}
		}
	}()

	return s.ListenAndServe("tcp", addr)
}

func (s *OpcServer) ListenAndServe(protocol string, port string) error {
	listener, err := net.Listen(protocol, port)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("OPC accept error", err)
			return err
		} else {
			log.Println("OPC client connected", conn.RemoteAddr())
		}

		go s.readMessages(conn)
	}
}

func (s *OpcServer) readMessages(conn net.Conn) {
	for {
		err, msg := s.readMessage(conn)
		if err != nil {
			log.Println("OPC client read error", conn.RemoteAddr(), err)
			break
		}

		s.Messages <- msg
	}

	conn.Close()
}

// Format: [channel, command, length high byte, length low byte, data...]
func (s *OpcServer) readMessage(conn net.Conn) (error, *OpcMessage) {
	header := make([]byte, 4)
	bytesRead, err := conn.Read(header)

	if err != nil {
		return err, nil
	} else if bytesRead != 4 {
		return errors.New("bad header"), nil
	}

	msg := &OpcMessage{
		Channel: header[0],
		Command: header[1],
		Length: binary.BigEndian.Uint16(header[2:4]),
	}

	msg.Data = make([]byte, msg.Length)
	bytesRead, err = conn.Read(msg.Data)

	if err != nil {
		return err, msg
	} else if bytesRead != int(msg.Length) {
		return errors.New("data length mismatch"), msg
	} else {
		return nil, msg
	}
}

type OpcMessage struct {
	Channel, Command uint8

	Length uint16
	Data []byte
}

func (m *OpcMessage) WritePixels(buffer [][]uint8) {
	for i, val := range m.Data {
		if len(buffer) == i / 3 {
			break
		}
		buffer[i / 3][i % 3] = val
	}
}

