package main

import (
	"bytes"
	"flag"

	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/opc"
	"github.com/brianewing/redshift/strip"

	"syscall/js"
	// "time"
)

// const leds = 256
// const leds = 256
// const leds = 10000

var leds = flag.Int("leds", 256, "number of LEDs")

var myAnimator *animator.Animator
var opcSession *opc.Session

var jsPixelBuffer js.Value

var stop = make(chan struct{})

func main() {
	flag.Parse()

	var effectSet effects.EffectSet

	if savedEffectsJSON := loadEffectsJSON(); len(savedEffectsJSON) > 0 {
		effectSet, _ = effects.UnmarshalJSON(savedEffectsJSON)
	}

	if len(effectSet) == 0 {
		rainbow := effects.NewRainbow()
		rainbow.Size = 200
		rainbow.Speed /= 3

		effectSet = effects.EffectSet{
			effects.EffectEnvelope{Effect: &effects.Clear{}},
			effects.EffectEnvelope{Effect: rainbow},
		}
	}

	myAnimator = &animator.Animator{
		Strip:   strip.New(*leds),
		Effects: effectSet,
	}

	// myAnimator.Init()

	// go myAnimator.Run(time.Second / 60)

	jsPixelBuffer = js.Global().Get("Uint8Array").New(*leds * 3)

	opcSession = opc.NewSession(myAnimator, jsOpcWriter{})

	js.Global().Set("RedshiftSend", js.FuncOf(ReceiveOpc))
	js.Global().Set("RedshiftStep", js.FuncOf(Step))
	js.Global().Set("RedshiftStop", js.FuncOf(Stop))

	// wait forever
	<-stop
}

func Step(this js.Value, args []js.Value) interface{} {
	myAnimator.Render()

	bytes := myAnimator.Strip.MarshalBytes()
	js.CopyBytesToJS(jsPixelBuffer, bytes)

	return jsPixelBuffer
}

func Stop(this js.Value, args []js.Value) interface{} {
	stop <- struct{}{}
	return nil
}

var receiveBuffer = make([]byte, 65536)

func ReceiveOpc(this js.Value, args []js.Value) interface{} {
	js.CopyBytesToGo(receiveBuffer, args[0])

	reader := bytes.NewReader(receiveBuffer)

	if msg, err := opc.ReadMessage(reader); err == nil {
		opcSession.Receive(msg)

		if msg.SystemExclusive.Command == opc.CmdSetEffectsJson {
			saveEffectsJSON(msg.SystemExclusive.Data)
		}
	} else {
		println("ReceiveOpc - Error -", err.Error())
	}

	return nil
}

type jsOpcWriter struct{}

func (_ jsOpcWriter) WriteOpc(msg opc.Message) error {
	bytes := msg.Bytes()
	jsMsgBuffer := js.Global().Get("Uint8Array").New(len(bytes))

	js.CopyBytesToJS(jsMsgBuffer, bytes)
	js.Global().Call("RedshiftReceive", jsMsgBuffer)

	return nil
}

func saveEffectsJSON(data []byte) {
	js.Global().Get("localStorage").Call("setItem", "ledplane.effectsJson", string(data))
}

func loadEffectsJSON() []byte {
	return []byte(js.Global().Get("localStorage").Call("getItem", "ledplane.effectsJson").String())
}
