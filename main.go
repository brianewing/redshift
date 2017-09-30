package main

import (
	"time"
	"redshift/server"
	"redshift/strip"
	"redshift/effects"
	"redshift/animator"
	"flag"
	"log"
	"io/ioutil"
)

var animationInterval = flag.Duration("animationInterval", 16 * time.Millisecond, "interval between animation frames")

var wsInterval = flag.Duration("wsInterval", 16 * time.Millisecond, "ws2811/2812(b) refresh interval")
var wsPin = flag.Int("wsPin", 0, "ws2811/2812(b) gpio pin")
var wsBrightness = flag.Int("wsBrightness", 50, "ws2811/2812(b) brightness")

var httpAddr = flag.String("httpAddr", "0.0.0.0:9191", "http service address")
var opcAddr = flag.String("opcAddr", "0.0.0.0:7890", "opc service address")

var numLeds = flag.Int("leds", 60, "number of leds")
var pathToEffectsJson = flag.String("effectsJson", "", "path to effects json")

func main() {
	flag.Parse()

	writeEffectsJson("effects.default.json", defaultEffects())

	ledStrip := strip.New(*numLeds)
	opcStrip := strip.New(ledStrip.Size)

	if *wsPin != 0 {
		go ledStrip.RunWs2811(*wsPin, *wsInterval, *wsBrightness)
	}

	animator := &animator.Animator{
		Strip: ledStrip,
		Effects: append(
			getEffects(),
			&effects.Buffer{Buffer: opcStrip.Buffer},
		),
	}

	go server.RunWebSocketServer(*httpAddr, ledStrip, animator.Effects)
	go server.RunOpcServer(*opcAddr, opcStrip)

	animator.Run(*animationInterval)
}

func getEffects() []effects.Effect {
	if path := *pathToEffectsJson; path != "" {
		if effects, err := loadEffectsJson(path); err == nil {
			return effects
		} else {
			log.Fatalln("Could not load effects json", err)
			return nil
		}
	} else {
		return defaultEffects()
	}
}

func loadEffectsJson(path string) ([]effects.Effect, error) {
	if bytes, err := ioutil.ReadFile(path); err == nil {
		return effects.UnmarshalJson(bytes)
	} else {
		return nil, err
	}
}

func writeEffectsJson(path string, effects_ []effects.Effect) error {
	bytes, err := effects.MarshalJson(effects_)
	if err != nil {
		log.Fatalln("Could not write effects json", "(marshall error)", err)
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

func defaultEffects() []effects.Effect {
	return []effects.Effect{
		&effects.Clear{},
		//&effects.RaceTestEffect{},
		//&effects.RandomEffect{},
		&effects.RainbowEffect{Size: 150, Speed: 1, Dynamic: true},
		//&effects.BlueEffect{},
		&effects.LarsonEffect{Color: []uint8{0,0,0}},
	}
}
