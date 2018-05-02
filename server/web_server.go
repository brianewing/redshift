package server

import (
	"encoding/json"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/osc"
	"github.com/brianewing/redshift/strip"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
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
	if r.URL.Path == "/" {
		s.serveWelcome(w, r)
	} else {
		s.serveWebSocket(w, r)
	}
}

func (s *webServer) serveWelcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	welcomeInfo := map[string]interface{}{
		"redshift": "0.1.0",
	}

	json.NewEncoder(w).Encode(welcomeInfo)
}

func (s *webServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	sc := &streamConnection{Conn: c}

	if err != nil {
		log.Print("WS websocket upgrade error:", err)
		return
	} else {
		log.Print("WS websocket client connected (", r.URL, ") [", c.RemoteAddr().String(), "]")
	}

	switch r.URL.Path {
	case "/s/strip":
		go s.streamStripBuffer(sc)
		go sc.readFps()
	case "/s/effects":
		go s.streamEffectsJson(sc)
		go sc.readFps()
	case "/s/osc":
		// messages are relayed as soon as they are received
		// so there's no need for client to specify fps via streamConnection (sc)
		go s.streamOscMessages(c)
	case "/strip":
		go s.receiveBuffer(c)
	case "/effects":
		go s.receiveEffects(c)
	}
}

func (s *webServer) receiveBuffer(c *websocket.Conn) {
	for {
		if _, msg, err := c.ReadMessage(); err != nil {
			log.Println("WS buffer read error", err)
			return
		} else {
			strip.UnserializeBufferBytes(s.buffer, msg)
		}
	}
}

func (s *webServer) receiveEffects(c *websocket.Conn) {
	for {
		if _, msg, err := c.ReadMessage(); err != nil {
			log.Println("WS effects read error:", err)
			return
		} else if effects, err := effects.UnmarshalJSON(msg); err == nil {
			log.Println("WS effects received:", string(msg))
			s.animator.SetEffects(effects)
		} else {
			log.Println("WS effects parse error:", err)
		}
	}
}

type streamConnection struct {
	requestedFps uint8 // 0-255
	fpsChanged   chan bool
	*websocket.Conn
}

func (sc *streamConnection) NextFrame() bool {
	for sc.requestedFps == 0 {
		<-sc.fpsChanged // block til it's set
		return true
	}

	time.Sleep(time.Second / time.Duration(sc.requestedFps))
	return true
}

func (sc *streamConnection) readFps() {
	sc.fpsChanged = make(chan bool)
	sc.SetReadLimit(1) // 1 byte at a time

	for {
		if _, msg, err := sc.ReadMessage(); err == nil {
			oldFps := sc.requestedFps
			if len(msg) > 0 {
				sc.requestedFps = uint8(msg[0])
			} else {
				sc.requestedFps = 0
			}

			if oldFps == 0 {
				sc.fpsChanged <- true
			}
		} else {
			log.Println("WS read fps error:", err)
			break
		}
	}
}

func (s *webServer) streamStripBuffer(sc *streamConnection) {
	for sc.NextFrame() {
		s.animator.Strip.Lock()
		msg := s.animator.Strip.SerializeBytes()
		s.animator.Strip.Unlock()

		if err := sc.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			log.Println("WS write error:", err)
			break
		}
	}
}

func (s *webServer) streamEffectsJson(sc *streamConnection) {
	for sc.NextFrame() {
		effectsJson, _ := json.Marshal(s.animator.Effects)

		if err := sc.WriteMessage(websocket.TextMessage, effectsJson); err != nil {
			log.Println("WS effects write error:", err)
			break
		}
	}
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
