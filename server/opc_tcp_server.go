package server

import (
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/strip"
	"net"
	"log"
)

type OpcServer struct {
	animator *animator.Animator
	buffer strip.Buffer
}

func RunOpcServer(addr string, animator *animator.Animator, buffer strip.Buffer) {
	s := &OpcServer{animator: animator, buffer: buffer}
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
	opcSession := &OpcSession{animator: s.animator, client: opcTcpWriter{Conn: conn}}
	for {
		msg, err := ReadOpcMessage(conn)
		if err != nil {
			log.Println("OPC client read error", conn.RemoteAddr(), err)
			break
		}

		if msg.Command == 0 {
			msg.WritePixels(s.buffer)
		} else {
			opcSession.Receive(msg)
		}
	}
	opcSession.Close()
	conn.Close()
}

type opcTcpWriter struct {
	net.Conn
}

func (w opcTcpWriter) WriteOpc(msg OpcMessage) (error) {
	_, err := w.Conn.Write(msg.Bytes())
	return err
}
