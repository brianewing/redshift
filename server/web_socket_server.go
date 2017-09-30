package server

import (
	"log"
	"time"
	"net/http"
	"github.com/gorilla/websocket"
	"redshift/strip"
	"redshift/effects"
)

type webSocketServer struct {
	strip *strip.LEDStrip
	effects []effects.Effect

	server *http.Server
	upgrader *websocket.Upgrader
	http.Handler
}

func RunWebSocketServer(addr string, strip *strip.LEDStrip, effects []effects.Effect) error {
	wss := &webSocketServer{
		strip: strip,
		effects: effects,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func (r *http.Request) bool { return true }, // ALLOW CROSS-ORIGIN REQUESTS
		},
	}
	wss.server = &http.Server{Addr: addr, Handler: wss}
	return wss.server.ListenAndServe()
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
			//go s.receiveBuffer(c)
		case "/effects":
			//go s.receiveEffects(c)
	}
}

type streamConnection struct {
	requestedFps uint8 // 0-255
	fpsChanged chan bool
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
			sc.requestedFps = uint8(msg[0])
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
		s.strip.Lock()
		msg := serializeStripBytes(s.strip)
		s.strip.Unlock()

		if err := sc.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			log.Println("WS write error:", err)
			break
		}
	}
}

func (s *webSocketServer) streamEffectsJson(sc *streamConnection) {
	for sc.NextFrame() {
		effectsJson, _ := effects.MarshalJson(s.effects)

		if err := sc.WriteMessage(websocket.TextMessage, effectsJson); err != nil {
			log.Println("WS effects write error:", err)
			break
		}
	}
}

func serializeStripBytes(strip *strip.LEDStrip) []byte {
	bytes := make([]byte, len(strip.Buffer) * 3)
	for i, led := range strip.Buffer {
		y := i * 3
		bytes[y] = led[0]
		bytes[y+1] = led[1]
		bytes[y+2] = led[2]
	}
	return bytes
}
