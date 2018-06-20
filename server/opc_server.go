package server

import (
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/opc"
	"github.com/brianewing/redshift/strip"
	"log"
	"net"
)

type opcServer struct {
	animator *animator.Animator
	buffer   strip.Buffer
}

func RunOpcServer(addr string, animator *animator.Animator, buffer strip.Buffer) {
	s := &opcServer{animator: animator, buffer: buffer}
	log.Fatalln("OPC", s.ListenAndServe("tcp", addr))
}

func (s *opcServer) ListenAndServe(protocol string, port string) error {
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

func (s *opcServer) readMessages(conn net.Conn) {
	session := &opc.Session{Animator: s.animator, Client: tcpOpcWriter{Conn: conn}}
	for {
		msg, err := opc.ReadMessage(conn)
		if err != nil {
			log.Println("OPC client read error", conn.RemoteAddr(), err)
			break
		}

		if msg.Command == 0 {
			msg.WritePixels(s.buffer)
		} else {
			session.Receive(msg)
		}
	}
	session.Close()
	conn.Close()
}

type tcpOpcWriter struct {
	net.Conn
}

func (w tcpOpcWriter) WriteOpc(msg opc.Message) error {
	_, err := w.Conn.Write(msg.Bytes())
	return err
}
