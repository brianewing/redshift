package server

import (
	"log"
	"time"
	"net/http"
	"encoding/json"

	"github.com/gorilla/websocket"
	"redshift/strip"
)

type webSocketServer struct {
	server *http.Server
	strip *strip.LEDStrip
	bufferInterval time.Duration

	upgrader *websocket.Upgrader
	http.Handler
}

func RunWebSocketServer(strip *strip.LEDStrip, bufferInterval time.Duration) error {
	wss := &webSocketServer{
		strip: strip,
		bufferInterval: bufferInterval,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func (r *http.Request) bool { return true },
		},
	}

	wss.server = &http.Server{Addr: "0.0.0.0:9191", Handler: wss}
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
	go s.streamStripBuffer(c)
	//defer c.Close()
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
		s.strip.Sync.Lock()
		msg, _ := serializeStrip(s.strip)
		s.strip.Sync.Unlock()
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("WS write error: ", err)
			break
		}
		time.Sleep(s.bufferInterval)
	}
}

type jsonFormat struct {
	Buffer [][]int `json:"buffer"`
}

func serializeStrip(strip *strip.LEDStrip) ([]byte, error) {
	return json.Marshal(&jsonFormat{strip.Buffer})
}

