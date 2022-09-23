package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/opc"
	"github.com/brianewing/redshift/strip"
	"github.com/gorilla/websocket"
)

type webServer struct {
	animator *animator.Animator
	buffer   strip.Buffer

	server   *http.Server
	upgrader *websocket.Upgrader

	http.Handler
}

func RunWebServer(addr string, animator *animator.Animator, buffer strip.Buffer) {
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
	} else if r.URL.Path == "/sse" {
		s.serveEventStream(w, r)
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

func (s *webServer) serveEventStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(200)

	interval := time.Second / 80

	var buf bytes.Buffer

	flusher := w.(http.Flusher)

	for {
		s.animator.Strip.Lock()
		data, _ := json.Marshal(s.animator.Strip.Buffer)
		s.animator.Strip.Unlock()

		buf.Truncate(0)
		buf.WriteString("data: ")
		buf.Write(data)
		buf.WriteString("\n\n")

		if _, err := w.Write(buf.Bytes()); err == nil {
			flusher.Flush()
		} else {
			log.Println("SSE client left", err)
			break
		}

		time.Sleep(interval)
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

	opcSession := opc.NewSession(s.animator, &websocketOpcWriter{Conn: c})
	s.readOpcMessages(c, opcSession)

	opcSession.Close()
}

func (s *webServer) readOpcMessages(c *websocket.Conn, opcSession *opc.Session) {
	for {
		if _, data, err := c.ReadMessage(); err != nil {
			log.Println("WS websocket opc read error:", err)
			break
		} else if msg, err := opc.ReadMessage(bytes.NewReader(data)); err != nil {
			log.Println("WS websocket opc parse error:", err)
		} else if err := opcSession.Receive(msg); err != nil {
			log.Println("WS websocket opc handle error:", err)
		}
	}
}

type websocketOpcWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (w *websocketOpcWriter) WriteOpc(msg opc.Message) error {
	w.Lock()
	defer w.Unlock()

	return w.WriteMessage(websocket.BinaryMessage, msg.Bytes())
}
