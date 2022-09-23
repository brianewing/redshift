// +build !js

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"

	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/midi"
	"github.com/brianewing/redshift/repl"
	"github.com/brianewing/redshift/server"
	"github.com/brianewing/redshift/strip"
)

var numLeds = flag.Int("leds", 100, "number of leds")
var scriptsDir = flag.String("scriptsDir", "scripts", "scripts directory relative to cwd")
var effectsDir = flag.String("effectsDir", "usereffects", "effect definitions directory relative to cwd")
var animationInterval = flag.Duration("animationInterval", time.Second/60, "interval between animation frames")

var pathToEffectsJson = flag.String("effectsJson", "", "path to effects json file")
var pathToEffectsYaml = flag.String("effectsYaml", "", "path to effects yaml file")

var wsInterval = flag.Duration("wsInterval", 16*time.Millisecond, "ws2811/2812(b) refresh interval")
var wsPin = flag.Int("wsPin", 0, "ws2811/2812(b) gpio pin")
var wsBrightness = flag.Int("wsBrightness", 50, "ws2811/2812(b) brightness")

var midiListDevices = flag.Bool("midiListDevices", false, "prints a list of available midi devices")
var midiDeviceId = flag.Int("midiDeviceId", 0, "midi device id")

var httpAddr = flag.String("httpAddr", "0.0.0.0:9191", "http service address")
var davAddr = flag.String("davAddr", "0.0.0.0:9292", "webdav (scripts) service address")
var effectsDavAddr = flag.String("effectsDavAddr", "0.0.0.0:9393", "webdav (effects) service address")
var opcAddr = flag.String("opcAddr", "0.0.0.0:7890", "opc service address")
var oscAddr = flag.String("oscAddr", "0.0.0.0:9494", "osc service address")

var logToFile = flag.Bool("log", false, "log to redshift.log")
var logDebug = flag.Bool("logDebug", false, "include source code filenames in log messages")

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	if *logDebug {
		log.SetFlags(log.Ltime | log.Lshortfile)
	}
}

func main() {
	if _, err := os.Stat(*scriptsDir); os.IsNotExist(err) {
		writePackedScripts(*scriptsDir)
	}

	if *midiListDevices {
		devices := midi.Devices()
		println("** MIDI Devices Available **")
		for i, device := range devices {
			println("  ", i, "-", device.Name)
		}
		println("")
		return
	}

	ledStrip := strip.New(*numLeds)

	opcBuffer := strip.NewBuffer(ledStrip.Size)
	wssBuffer := strip.NewBuffer(ledStrip.Size)

	animator := &animator.Animator{
		Strip:   ledStrip,
		Effects: getEffects(),
		PostEffects: effects.EffectSet{
			effects.EffectEnvelope{Effect: effects.NewBlendFromBuffer(opcBuffer)},
			effects.EffectEnvelope{Effect: effects.NewBlendFromBuffer(wssBuffer)},
		},
	}

	if *wsPin != 0 {
		go ledStrip.RunWs2811(*wsPin, *wsInterval, *wsBrightness)
	}

	go server.RunWebServer(*httpAddr, animator, wssBuffer)
	go server.RunWebDavServer(*davAddr, *scriptsDir, true)
	go server.RunWebDavServer(*effectsDavAddr, *effectsDir, false)
	go server.RunOpcServer(*opcAddr, animator, opcBuffer)
	go server.RunOscServer(*oscAddr)

	// if *opcForwardAddr != "" {
	// 	go opc.RunForwarder(*opcForwardAddr, ledStrip)
	// }

	go repl.Run(animator, os.Stdin, "> ")
	go repl.RunReplServer(":9999", animator)

	go cleanupOnCtrlC(animator)

	animator.Run(*animationInterval)
}

func cleanupOnCtrlC(animator *animator.Animator) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	animator.Finish()
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
