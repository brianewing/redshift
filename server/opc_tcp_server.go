package server

import (
	"github.com/brianewing/redshift/strip"
	"net"
	"log"
)

type OpcServer struct {
	Messages chan OpcMessage
}

func RunOpcServer(addr string, buffer strip.Buffer) {
	s := &OpcServer{Messages: make(chan OpcMessage)}
	go func() {
		for {
			if msg := <-s.Messages; msg.Command == 0 {
				msg.WritePixels(buffer)
			}
		}
	}()
	log.Fatalln("OPC", s.ListenAndServe("tcp", addr))
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
		msg, err := ReadOpcMessage(conn)
		if err != nil {
			log.Println("OPC client read error", conn.RemoteAddr(), err)
			break
		}
		s.Messages <- msg
	}
	conn.Close()
}
