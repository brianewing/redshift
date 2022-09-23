package server

import "github.com/hypebeast/go-osc/osc"
import redshiftOsc "github.com/brianewing/redshift/osc"
import "log"

func RunOscServer(addr string) {
	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("*", func(msg *osc.Message) {
		redshiftOsc.ReceiveMessage(redshiftOsc.OscMessage{
			Address: msg.Address,
			Arguments: msg.Arguments,
		})
	})

	server := &osc.Server{
		Addr: addr,
		Dispatcher: dispatcher,
	}

	log.Fatalln(server.ListenAndServe())
}
