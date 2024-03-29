package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

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
	} else if r.URL.Path == "/opc-stream" {
		s.serveOpcStream(w, r)
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

func (s *webServer) serveOpcStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/opc-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("LedPlane-Strip-Size", strconv.Itoa(s.animator.Strip.Size))

	w.WriteHeader(200)

	fps := 60

	sess := opc.NewSession(s.animator, &opc.IOWriter{Writer: w})

	sess.Receive(opc.Message{
		Command: 255,
		Channel: 0,
		SystemExclusive: opc.SystemExclusive{
			Command: opc.CmdOpenStream,
			Data: []byte("strip"),
		},
	})

	sess.Receive(opc.Message{
		Command: 255,
		Channel: 0,
		SystemExclusive: opc.SystemExclusive{
			Command: opc.CmdSetStreamFps,
			Data: []byte{byte(fps)},
		},
	})

	<-r.Context().Done()
	sess.Close()

	return
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
