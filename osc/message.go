package osc

import "time"

type OscMessage struct {
	Address   string // uri
	Arguments []interface{}
	Timestamp time.Time
}
