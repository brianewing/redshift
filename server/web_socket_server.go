package server

import (
	"log"
	"time"
	"net/http"
	"encoding/json"

	"github.com/gorilla/websocket"
	"redshift/strip"
	"redshift/effects"
	"sync"
)

type webSocketServer struct {
	server *http.Server
	strip *strip.LEDStrip
	effects []effects.Effect
	bufferInterval time.Duration

	writeMutex sync.Mutex

	upgrader *websocket.Upgrader
	http.Handler
}

func RunWebSocketServer(addr string, strip *strip.LEDStrip, effects []effects.Effect, bufferInterval time.Duration) error {
	wss := &webSocketServer{
		strip: strip,
		effects: effects,
		bufferInterval: bufferInterval,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func (r *http.Request) bool { return true }, // ALLOW CROSS-ORIGIN REQUESTS
		},
	}

	wss.server = &http.Server{Addr: addr, Handler: wss}
	return wss.server.ListenAndServe()
}

func (s *webSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("WS upgrade error: ", err)
		return
	} else {
		log.Print("WS client connected ", c.RemoteAddr().String())
	}

	go s.readMessages(c)

	println(r.URL.Path)

	switch r.URL.Path {
		case "/strip": go s.streamStripBuffer(c)
		case "/effects": go s.streamEffectsJson(c)
	}
}

func (s *webSocketServer) readMessages(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("WS read error: ", err)
			break
		} else {
			log.Println("WS got message: ", message)
		}
	}
}

func (s *webSocketServer) streamStripBuffer(c *websocket.Conn) {
	for {
		s.strip.Lock()
		msg := serializeStripBytes(s.strip)
		s.strip.Unlock()

		s.writeMutex.Lock()
		err := c.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			log.Println("WS write error: ", err)
			break
		}
		s.writeMutex.Unlock()

		time.Sleep(s.bufferInterval)
	}
}

func (s *webSocketServer) streamEffectsJson(c *websocket.Conn) {
	for {
		effectsJson, _ := effects.MarshalJson(s.effects)
		s.writeMutex.Lock()
		err := c.WriteMessage(websocket.TextMessage, effectsJson)
		s.writeMutex.Unlock()
		if err != nil {
			log.Println("WS", "effects", "write error", err)
			break
		}
		time.Sleep(s.bufferInterval)
	}
}

type jsonFormat struct {
	Buffer [][]uint8 `json:"buffer"`
}

func serializeStripJson(strip *strip.LEDStrip) ([]byte, error) {
	return json.Marshal(&jsonFormat{strip.Buffer})
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

