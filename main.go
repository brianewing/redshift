package main

import (
	"time"
	"redshift/server"
	"redshift/strip"
	"redshift/effects"
	"redshift/animator"
	"flag"
	"log"
	"github.com/luci/go-render/render"
)

const ANIMATION_INTERVAL = 16 * time.Millisecond
const WSS_BUFFER_INTERVAL = 16 * time.Millisecond

var numLeds = flag.Int("leds", 60, "number of leds")
var httpAddr = flag.String("httpAddr", "0.0.0.0:9191", "http service address")
var opcAddr = flag.String("opcAddr", "0.0.0.0:7890", "opc service address")

func main() {
	flag.Parse()

	ledStrip := strip.New(*numLeds)
	opcStrip := strip.New(ledStrip.Size)

	animator := &animator.Animator{
		Strip: ledStrip,
		Effects: []effects.Effect{
			&effects.Clear{},
			&effects.Buffer{Buffer: opcStrip.Buffer},
			//&effects.RaceTestEffect{},
			//&effects.RandomEffect{},
			&effects.RainbowEffect{Size: 150, Speed: 1, Dynamic: true},
			//&effects.BlueEffect{},
			&effects.LarsonEffect{Color: []uint8{0,0,0}},
		},
	}

	//log.Println(render.Render(animator.Effects))
	json, err := effects.MarshalJson(animator.Effects)
	log.Println(string(json), err)

	effects2, err := effects.UnmarshalJson(json)
	log.Println(render.Render(effects2), err)

	go server.RunWebSocketServer(*httpAddr, ledStrip, animator.Effects, WSS_BUFFER_INTERVAL)
	go server.RunOpcServer(*opcAddr, opcStrip)

	animator.Run(ANIMATION_INTERVAL)
}