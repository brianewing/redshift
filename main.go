package main

import (
	"time"
	"redshift/server"
	"redshift/strip"
	"redshift/effects"
	"redshift/animator"
)

const LEDS = 60
const ANIMATION_INTERVAL = 30 * time.Millisecond

func main() {
	ledStrip := strip.New(LEDS)
	ledStrip.Clear()

	opcStrip := strip.New(LEDS)
	opcStrip.Clear()

	//addr := flag.String("httpAddr", "localhost:9191", "http service address")
	//flag.Parse()

	go server.RunWebSocketServer(ledStrip)
	go server.RunOpcServer(ledStrip, opcStrip.Buffer)

	animator := &animator.Animator{
		Strip: ledStrip,
		Effects: []effects.Effect{
			&effects.Clear{},
			&effects.Buffer{Buffer: opcStrip.Buffer},
			//&effects.RaceTestEffect{},
			//&effects.RandomEffect{},
			&effects.RainbowEffect{Size: 150, Speed: 1, Dynamic: true},
			//&effects.BlueEffect{},
			&effects.LarsonEffect{Color: []int{0,0,0}},
		},
	}

	animator.Run(ANIMATION_INTERVAL)
}