package osc

import (
	"log"
	"sync"
	"time"
)

// This is a bit of a hack, could be improved with better architecture
// I hope to come back and refactor this

// The spool receives messages from the OscServer in the server package
// and fans them out to listeners (e.g. OscControl's in redshift/effects)

var spool struct {
	listeners []chan OscMessage
	sync.Mutex
}

func StreamMessages() (msgs chan OscMessage, done chan struct{}) {
	msgs = make(chan OscMessage)
	done = make(chan struct{})

	spool.Lock()
	spool.listeners = append(spool.listeners, msgs)
	spool.Unlock()

	go removeWhenDone(msgs, done)

	return
}

func ReceiveMessage(msg OscMessage) {
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	spool.Lock()
	for _, c := range spool.listeners {
		c <- msg
	}
	spool.Unlock()
}

func removeWhenDone(msgs chan OscMessage, done chan struct{}) {
	<-done
	spool.Lock()
	for i, c := range spool.listeners {
		if c == msgs {
			spool.listeners = append(spool.listeners[:i], spool.listeners[i+1:]...)
			close(c)
			break
		}
	}
	spool.Unlock()
}

func init() {
	go debugOscMessages()
}

func debugOscMessages() {
	debugStream, _ := StreamMessages()
	for msg := range debugStream {
		log.Println("Incoming OSC message:", msg.Address, msg.Arguments)
	}
}
