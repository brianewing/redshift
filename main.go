package main

import (
	"time"
	"redshift/server"
	"redshift/strip"
	"redshift/effects"
	"redshift/animator"
	"flag"
)

const ANIMATION_INTERVAL = 30 * time.Millisecond
const WSS_BUFFER_INTERVAL = 30 * time.Millisecond

var numLeds = flag.Int("leds", 60, "number of leds")
var httpAddr = flag.String("httpAddr", "0.0.0.0:9191", "http service address")
var opcAddr = flag.String("opcAddr", "0.0.0.0:7890", "opc service address")

func main() {
	flag.Parse()

	ledStrip := strip.New(*numLeds)
	ledStrip.Clear()

	opcStrip := strip.New(ledStrip.Size)
	opcStrip.Clear()

	go server.RunWebSocketServer(*httpAddr, ledStrip, WSS_BUFFER_INTERVAL)
	go server.RunOpcServer(*opcAddr, opcStrip)

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