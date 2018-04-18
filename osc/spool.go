package osc

import (
	"log"
	"sync"
)

// This is a bit of a hack, could be improved with better architecture
// I hope to come back and refactor this some day

// The spool receives messages from the OscServer in the server package
// and fans them out to listeners (e.g. OscControl's in redshift/effects)

var oscSpool struct {
	listeners []chan OscMessage
	sync.Mutex
}

func StreamMessages() chan OscMessage {
	oscSpool.Lock()

	newChan := make(chan OscMessage)
	oscSpool.listeners = append(oscSpool.listeners, newChan)
	log.Println("new listeners:", oscSpool.listeners)

	oscSpool.Unlock()
	return newChan
}

func PushToListeners(msg OscMessage) {
	oscSpool.Lock()
	for _, c := range oscSpool.listeners {
		c <- msg
	}
	oscSpool.Unlock()
}

func init() {
	debugStream := StreamMessages()
	go func() {
		for msg := range debugStream {
			log.Println("Incoming OSC message:", msg.Address, msg.Arguments)
		}
	}()
}
