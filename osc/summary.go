package osc

// osc address => last msg
var summary = make(map[string]OscMessage)

// Returns a summary map all the messages that have been received
// Keys represent an osc address and values are the latest
// message received for that address
func Summary() map[string]OscMessage {
	return summary
}

// Resets the summary by replacing it with an empty map
func ClearSummary() {
	summary = make(map[string]OscMessage)
}

func readSummary() {
	msgs, _ := StreamMessages()

	for msg := range msgs {
		summary[msg.Address] = msg
	}
}

func init() {
	go readSummary()
}
