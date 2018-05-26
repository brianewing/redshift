package server

import (
	"bytes"
	"encoding/json"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/osc"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type webServer struct {
	animator *animator.Animator
	buffer   [][]uint8

	server   *http.Server
	upgrader *websocket.Upgrader

	http.Handler
}

func RunWebServer(addr string, animator *animator.Animator, buffer [][]uint8) {
	s := &webServer{
		animator: animator,
		buffer:   buffer,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // ALLOW CROSS-ORIGIN REQUESTS
		},
	}
	s.server = &http.Server{Addr: addr, Handler: s}
	log.Fatalln("WS", s.server.ListenAndServe())
}

func (s *webServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/welcome" {
		s.serveWelcome(w, r)
	} else {
		s.serveWebSocket(w, r)
	}
}

func (s *webServer) serveWelcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.welcomeInfo())
}

func (s *webServer) welcomeInfo() map[string]interface{} {
	return map[string]interface{}{
		"redshift": "0.1.0",
	}
}

func (s *webServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("WS websocket upgrade error:", err)
		return
	} else {
		log.Print("WS websocket client connected (", r.URL, ") [", c.RemoteAddr().String(), "]")
	}

	opcSession := &OpcSession{animator: s.animator, client: websocketOpcWriter{Conn: c}}
	s.readOpcMessages(c, opcSession)

	opcSession.Close()
}

func (s *webServer) readOpcMessages(c *websocket.Conn, opcSession *OpcSession) {
	for {
		if _, data, err := c.ReadMessage(); err != nil {
			log.Println("WS websocket opc read error:", err)
			break
		} else if msg, err := ReadOpcMessage(bytes.NewReader(data)); err != nil {
			log.Println("WS websocket opc parse error:", err)
		} else if err := opcSession.Receive(msg); err != nil {
			log.Println("WS websocket opc handle error:", err)
		}
	}
}

type websocketOpcWriter struct {
	*websocket.Conn
}

func (w websocketOpcWriter) WriteOpc(msg OpcMessage) error {
	return w.WriteMessage(websocket.BinaryMessage, msg.Bytes())
}

func (s *webServer) streamOscMessages(c *websocket.Conn) {
	oscMessages, stop := osc.StreamMessages()
	for msg := range oscMessages {
		msgJson, _ := json.Marshal(msg)
		if err := c.WriteMessage(websocket.TextMessage, msgJson); err != nil {
			log.Println("WS osc write error:", err)
			break
		}
	}
	stop <- struct{}{}
}
