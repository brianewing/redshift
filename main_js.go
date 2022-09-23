package main

import (
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/strip"
	"github.com/brianewing/redshift/effects"

	"log"
	"math/rand"
	"time"

	"syscall/js"
)

const numLEDs = 900

var world struct{
	*animator.Animator
	*strip.LEDStrip
}

func init() {
	rand.Seed(time.Now().UnixNano())

	world.Strip = strip.New(numLEDs)
	world.Buffer = strip.NewBuffer(numLEDs)

	world.Animator = &animator.Animator{
		Strip: world.Strip,
		Effects: effects.EffectSet{
			effects.EffectEnvelope{
				Effect: effects.NewRainbow(),
			},
		},
	}

	js.Global().Set("LedPlaneStep", js.FuncOf(Step))
}

func main() {
	log.Println("Hello from LedPlane WASM!")
}

func Step(this js.Value, args []js.Value) interface{} {
	world.Animator.Render()
	return world.Strip.Buffer
}

