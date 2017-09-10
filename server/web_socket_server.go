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

	upgrader *websocket.Upgrader
	http.Handler
}

func Run(strip *strip.LEDStrip) {
	wss := &webSocketServer{
		strip: strip,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func (r *http.Request) bool { return true },
		},
	}

	wss.server = &http.Server{Addr: "localhost:9191", Handler: wss}
	wss.server.ListenAndServe()
}

func (s *webSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("WS upgrade error: ", err)
		return
	}

	go s.readMessages(c)
	go s.streamStripBuffer(c)
	//defer c.Close()
}

func (s *webSocketServer) readMessages(c *websocket.Conn) {
	for {
		log.Println("Waiting to read...")
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error: ", err)
			break
		} else {
			log.Println("Got message: ", message)
		}
	}
}

func (s *webSocketServer) streamStripBuffer(c *websocket.Conn) {
	for {
		s.strip.Sync.Lock()
		msg, _ := serializeStrip(s.strip)
		err := c.WriteMessage(websocket.TextMessage, msg)
		s.strip.Sync.Unlock()
		if err != nil {
			log.Println("Write error: ", err)
			break
		}
		time.Sleep(30 * time.Millisecond)
	}
}

type jsonFormat struct {
	Buffer *[][]int `json:"buffer"`
}

func serializeStrip(strip *strip.LEDStrip) ([]byte, error) {
	return json.Marshal(&jsonFormat{&strip.Buffer})
}

