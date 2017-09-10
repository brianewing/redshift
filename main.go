package main

import (
	"time"
	"redshift/server"
	"redshift/strip"
	"redshift/effects"
)

const LEDS = 60

func main() {
	//addr := flag.String("httpAddr", "localhost:9191", "http service address")
	//flag.Parse()

	ledStrip := strip.New(LEDS)
	go server.Run(ledStrip)

	animator := &Animator{
		Strip: ledStrip,
		Effects: []effects.Effect{
			&effects.Clear{},
			//&effects.RaceTestEffect{},
			&effects.RandomEffect{},
			&effects.BlueEffect{},
			&effects.LarsonEffect{Color: []int{0,0,0}},
		},
		Interval: 30 * time.Millisecond,
	}

	animator.Run()
}