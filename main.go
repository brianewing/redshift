package main

import (
	"flag"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/server"
	"github.com/brianewing/redshift/strip"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

var numLeds = flag.Int("leds", 30, "number of leds")
var scriptsDir = flag.String("scriptsDir", "scripts", "scripts directory relative to cwd")
var animationInterval = flag.Duration("animationInterval", 16*time.Millisecond, "interval between animation frames")

var pathToEffectsJson = flag.String("effectsJson", "", "path to effects json file")
var pathToEffectsYaml = flag.String("effectsYaml", "", "path to effects yaml file")

var wsInterval = flag.Duration("wsInterval", 16*time.Millisecond, "ws2811/2812(b) refresh interval")
var wsPin = flag.Int("wsPin", 0, "ws2811/2812(b) gpio pin")
var wsBrightness = flag.Int("wsBrightness", 50, "ws2811/2812(b) brightness")

var httpAddr = flag.String("httpAddr", "0.0.0.0:9191", "http service address")
var opcAddr = flag.String("opcAddr", "0.0.0.0:7890", "opc service address")
var davAddr = flag.String("davAddr", "0.0.0.0:9292", "webdav (scripts) service address")

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if _, err := os.Stat(*scriptsDir); os.IsNotExist(err) {
		writePackedScripts(*scriptsDir)
	}

	ledStrip := strip.New(*numLeds)

	opcBuffer := strip.NewBuffer(ledStrip.Size)
	wssBuffer := strip.NewBuffer(ledStrip.Size)

	animator := &animator.Animator{
		Strip:   ledStrip,
		Effects: getEffects(),
		PostEffects: []effects.Effect{
			&effects.Buffer{Buffer: opcBuffer},
			&effects.Buffer{Buffer: wssBuffer},
		},
	}

	if *wsPin != 0 {
		go ledStrip.RunWs2811(*wsPin, *wsInterval, *wsBrightness)
	}

	go server.RunWebSocketServer(*httpAddr, ledStrip, wssBuffer, &animator.Effects)
	go server.RunWebDavServer(*davAddr, *scriptsDir)
	go server.RunOpcServer(*opcAddr, opcBuffer)

	animator.Run(*animationInterval)
}

//go:generate go-bindata --prefix skel/ skel/...
func writePackedScripts(dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	scripts, _ := AssetDir("scripts")
	for _, name := range scripts {
		data, _ := Asset(path.Join("scripts", name))
		if err := ioutil.WriteFile(path.Join(dest, name), data, 0744); err != nil {
			return err
		}
	}
	return nil
}
