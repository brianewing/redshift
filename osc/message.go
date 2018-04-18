package osc

type OscMessage struct {
	Address string // uri
	Arguments []interface{}
}
