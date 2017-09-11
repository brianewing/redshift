package main

import (
	"time"
	"redshift/server"
	"redshift/strip"
	"redshift/effects"
	"redshift/animator"
)

const LEDS = 60

func main() {
	//addr := flag.String("httpAddr", "localhost:9191", "http service address")
	//flag.Parse()

	ledStrip := strip.New(LEDS)
	ledStrip.Clear()

	onUpdate := make(chan int)

	go server.RunWebSocketServer(ledStrip)
	go server.RunOpcServer(ledStrip, onUpdate)

	animator := &animator.Animator{
		Strip: ledStrip,
		Effects: []effects.Effect{
			//&effects.Clear{},
			//&effects.RaceTestEffect{},
			//&effects.RandomEffect{},
			&effects.BlueEffect{},
			&effects.LarsonEffect{Color: []int{0,0,0}},
		},
		Interval: 30 * time.Millisecond,
	}

	//animator.Run()

	for {
		<- onUpdate
		//fmt.Println("Update")
		animator.Render()
	}
}