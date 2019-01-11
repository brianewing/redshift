package osc

import "sync"

var summaryMutex sync.Mutex
var summary = make(map[string]OscMessage) // osc address => last msg

// returns a summary of the messages that have been received so far as a map,
// where keys represent an osc address and values are the latest message received for that address
func Summary() map[string]OscMessage {
  summaryMutex.Lock()
  _copy := make(map[string]OscMessage, len(summary))
  for addr, msg := range summary {
    _copy[addr] = msg
  }
  summaryMutex.Unlock()
	return _copy
}

func ClearSummary() {
  summaryMutex.Lock()
	summary = make(map[string]OscMessage)
  summaryMutex.Unlock()
}

// streams messages the spool contained in this package
// and sets summary[msg.Address] = msg for each one
func createSummary() {
	msgs, _ := StreamMessages()
	for msg := range msgs {
    summaryMutex.Lock()
		summary[msg.Address] = msg
    summaryMutex.Unlock()
	}
}

func init() {
	go createSummary()
}
