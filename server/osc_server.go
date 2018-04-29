package server

import "github.com/hypebeast/go-osc/osc"
import redshiftOsc "github.com/brianewing/redshift/osc"
import "log"

func RunOscServer(addr string) {
	server := &osc.Server{Addr: addr}

	server.Handle("*", func(msg *osc.Message) {
		redshiftOsc.ReceiveMessage(redshiftOsc.OscMessage{
			Address: msg.Address,
			Arguments: msg.Arguments,
		})
	})

	log.Fatalln(server.ListenAndServe())
}
