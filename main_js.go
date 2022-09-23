package main

import (
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/opc"
	"github.com/brianewing/redshift/strip"
	"github.com/brianewing/redshift/effects"

	"log"
	"math/rand"
	"time"

	"syscall/js"
)

const numLEDs = 900

var World struct{
	*animator.Animator
	*strip.LEDStrip

	*opc.Session

	receiver js.Value
}

func init() { rand.Seed(time.Now().UnixNano()) }

func init() {
	World.Strip = strip.New(numLEDs)
	World.Buffer = strip.NewBuffer(numLEDs)

	World.Animator = &animator.Animator{
		Strip: World.Strip,
		Effects: effects.EffectSet{
			effects.EffectEnvelope{
				Effect: effects.NewRainbow(),
			},
		},
	}

	World.Session = opc.NewSession(World.Animator, jsOpcWriter{})

	js.Global().Set("LedPlaneSend", js.FuncOf(Send))

	World.receiver = js.Global().Get("LedPlaneReceive")
}

func main() {
	log.Println("Hello from LedPlane WASM!")
}

func Send(this js.Value, args []js.Value) interface{} {
	World.Session.ReceiveMessage()
	World.Animator.Render()
	return World.Strip.Buffer
}
