package effects

import (
	"errors"
	"github.com/brianewing/redshift/midi"
	"github.com/brianewing/redshift/osc"
	"log"
	"reflect"
)

type Control interface {
	Apply(effect interface{})
}

type ControlEnvelope struct {
	Control
	Controls ControlSet // recursively controllable
}

func (e *ControlEnvelope) Init() {
	if initable, ok := e.Control.(Initable); ok {
		initable.Init()
	}
	e.Controls.Init()
}

func (e *ControlEnvelope) Destroy() {
	if destroyable, ok := e.Control.(Destroyable); ok {
		destroyable.Destroy()
	}
	e.Controls.Destroy()
}

func (e *ControlEnvelope) Apply(effect interface{}) {
	e.Controls.Apply(e.Control) // meta controls
	e.Control.Apply(effect)
}

type ControlSet []ControlEnvelope

func (set ControlSet) Apply(effect interface{}) {
	for _, control := range set {
		control.Apply(effect)
	}
}

func (set ControlSet) Init() {
	for _, envelope := range set {
		envelope.Init()
	}
}

func (set ControlSet) Destroy() {
	for _, envelope := range set {
		envelope.Destroy()
	}
}

/*
 * Fixed Value Control
 */

type FixedValueControl struct {
	Field string
	Value interface{}
}

func (c *FixedValueControl) Apply(effect interface{}) {
	if c.Value == nil {
		return
	}
	if field, err := getField(effect, c.Field); err == nil {
		if err := setValue(field, c.Value); err != nil {
			log.Println("FixedValueControl", err)
		}
	}
}

/*
 * Tween Control
 */

type TweenControl struct {
	Field, Function string
	Min, Max, Speed float64
}

func (c *TweenControl) getFunction() TimingFunction {
	switch c.Function {
	case "triangle":
		return OscillateBetween
	case "sawtooth":
		return CycleBetween
	case "sin":
	}
	return SmoothOscillateBetween
}

func (c *TweenControl) Apply(effect interface{}) {
	if field, err := getField(effect, c.Field); err == nil {
		fn := c.getFunction()
		val := round(fn(c.Min, c.Max, c.Speed))

		setValue(field, val)
	}
}

/*
 * Midi Control
 */

type MidiControl struct {
	Field string

	Device        int
	Status, Data1 int64

	Min, Max float64

	value float64 // latched from most recent midi msg

	midiMsgs chan midi.MidiMessage
}

func (c *MidiControl) Init() {
	if devices := midi.Devices(); len(devices) > c.Device {
		device := devices[c.Device]
		c.midiMsgs = midi.StreamMessages(device)
		go c.readValues()
	} else {
		log.Println("MidiControl", "device not found", "| id:", c.Device)
	}
}

func (c *MidiControl) Destroy() {
	close(c.midiMsgs)
}

func (c *MidiControl) Apply(effect interface{}) {
	(&FixedValueControl{Field: c.Field, Value: c.value}).Apply(effect)
}

func (c *MidiControl) readValues() {
	for msg := range c.midiMsgs {
		if msg.Status == c.Status && msg.Data1 == c.Data1 {
			c.value = c.scaleValue(msg.Data2)
		}
	}
}

func (c *MidiControl) scaleValue(val int64) float64 {
	v := float64(val)
	return (v/127.0)*(c.Max-c.Min) + c.Min
}

/*
 * Osc Control
 */

type OscControl struct{
	Address string
	stream chan osc.OscMessage

	FixedValueControl
}

func (c *OscControl) Init() {
	c.stream = osc.StreamMessages()
	go func() {
		for msg := range c.stream {
			if msg.Address == c.Address && len(msg.Arguments) >= 2 {
				c.Value = int(msg.Arguments[1].(int32) * 2)
				log.Println("Woot msg: ", c.Value)
			}
		}
	}()
}

/*
 * Null Control
 */

type NullControl struct{}

func (c NullControl) Apply(interface{}) {}

/*
 * Construction
 */

func ControlByName(name string) Control {
	switch name {
	case "FixedValueControl":
		return &FixedValueControl{}
	case "TweenControl":
		return &TweenControl{}
	case "MidiControl":
		return &MidiControl{}
	case "OscControl":
		return &OscControl{}
	}
	return NullControl{}
}

/*
 * Reflection functions
 */

func getField(effect interface{}, name string) (reflect.Value, error) {
	field := reflect.ValueOf(effect).Elem().FieldByName(name)
	if field.IsValid() {
		return field, nil
	} else {
		return field, errors.New("field not found")
	}
}

func setValue(field reflect.Value, newVal interface{}) error {
	if field.Type() == reflect.TypeOf(newVal) {
		switch field.Interface().(type) {
		case int:
			field.SetInt(int64(newVal.(int)))
		case uint:
			field.SetUint(uint64(newVal.(int)))
		case float64:
			field.SetFloat(float64(newVal.(float64)))
		default:
			return errors.New("can't set field (unknown type)")
		}
		return nil
	} else {
		log.Println("set err", field.Type(), reflect.TypeOf(newVal))
		return errors.New("type mismatch")
	}
}
