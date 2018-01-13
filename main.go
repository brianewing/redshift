package main

import (
	"flag"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/midi"
	"github.com/brianewing/redshift/server"
	"github.com/brianewing/redshift/strip"
	"io/ioutil"
	"log"
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

var midiDeviceId = flag.Int("midiDeviceId", 0, "midi device id")

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

	rainbowEffect := &effects.RainbowEffect{Size: 100, Speed: 1}
	larsonEffect := &effects.LarsonEffect{Color: []uint8{255, 255, 255}}
	brightnessEffect := &effects.Brightness{Brightness: 255}

	adjustParametersFromMidiEvent := func(event midi.MidiMessage) {
		switch event.Status {
		case 176:
			switch event.Data1 {
			case 1:
				rainbowEffect.Size = uint(event.Data2)
			case 2:
				rainbowEffect.Speed = float64(event.Data2 / 6)
			case 4:
				larsonEffect.Position = int(int(event.Data2) / ledStrip.Size)
			case 95:
				log.Println("Set brightness")
				brightnessEffect.Brightness = uint8(event.Data2 * 2)
			}
		}
	}

	devices := midi.Devices()

	for i, device := range devices {
		log.Println("MIDI Device", i, device)
	}

	if *midiDeviceId != 0 {
		device := devices[*midiDeviceId]
		midiEventsChan := midi.StreamMessages(device)

		go func() {
			log.Println("MIDI start reading")
			for midiEvent := range midiEventsChan {
				log.Println("MIDI Event", midiEvent)
				adjustParametersFromMidiEvent(midiEvent)
			}
			log.Println("MIDI finished reading")
		}()
	}

	animator := &animator.Animator{
		Strip:   ledStrip,
		Effects: []effects.Effect{
			&effects.Clear{},
			rainbowEffect,
			larsonEffect,
			brightnessEffect,
			//&effects.External{Program: "scripts/example.py"},
		},
		PostEffects: []effects.Effect{
			&effects.Buffer{Buffer: opcBuffer},
			&effects.Buffer{Buffer: wssBuffer},
		},
	}

	if *wsPin != 0 {
		go ledStrip.RunWs2811(*wsPin, *wsInterval, *wsBrightness)
	}

	go server.RunWebSocketServer(*httpAddr, animator, wssBuffer)
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
