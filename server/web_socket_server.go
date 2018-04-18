package server

import (
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type webSocketServer struct {
	animator *animator.Animator
	buffer   [][]uint8

	server   *http.Server
	upgrader *websocket.Upgrader
	http.Handler
}

func RunWebSocketServer(addr string, animator *animator.Animator, buffer [][]uint8) {
	wss := &webSocketServer{
		animator: animator,
		buffer:   buffer,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // ALLOW CROSS-ORIGIN REQUESTS
		},
	}
	wss.server = &http.Server{Addr: addr, Handler: wss}
	log.Fatalln("WS", wss.server.ListenAndServe())
}

func (s *webSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	sc := &streamConnection{Conn: c}

	if err != nil {
		log.Print("WS upgrade error:", err)
		return
	} else {
		log.Print("WS client connected (", r.URL, ") [", c.RemoteAddr().String(), "]")
	}

	switch r.URL.Path {
	case "/s/strip":
		go s.streamStripBuffer(sc)
		go sc.readFps()
	case "/s/effects":
		go s.streamEffectsJson(sc)
		go sc.readFps()
	case "/strip":
		go s.receiveBuffer(c)
	case "/effects":
		go s.receiveEffects(c)
	}
}

func (s *webSocketServer) receiveBuffer(c *websocket.Conn) {
	for {
		if _, msg, err := c.ReadMessage(); err != nil {
			log.Println("WS buffer read error", err)
			return
		} else {
			strip.UnserializeBufferBytes(s.buffer, msg)
		}
	}
}

func (s *webSocketServer) receiveEffects(c *websocket.Conn) {
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

func (s *webSocketServer) streamStripBuffer(sc *streamConnection) {
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

func (s *webSocketServer) streamEffectsJson(sc *streamConnection) {
	for sc.NextFrame() {
		effectsJson, _ := effects.MarshalJSON(s.animator.Effects)

		if err := sc.WriteMessage(websocket.TextMessage, effectsJson); err != nil {
			log.Println("WS effects write error:", err)
			break
		}
	}
}
