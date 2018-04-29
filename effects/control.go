package effects

import (
	"errors"
	"github.com/brianewing/redshift/midi"
	"github.com/brianewing/redshift/osc"
	"github.com/robertkrimen/otto"
	"log"
	"reflect"
	"strconv"
	"strings"
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
 * Base Control
 */

type BaseControl struct {
	Field     string
	LastError string

	Initial interface{}
	value   interface{}

	Transform string
	vm        *otto.Otto
}

func (c *BaseControl) Init() {
	c.value = c.Initial
}

func (c *BaseControl) Apply(effect interface{}) {
	if c.value == nil {
		return
	}
	if field, err := getField(effect, c.Field); err == nil {
		val := c.transformValue()
		if err := setValue(field, val); err != nil {
			c.setError(err)
		}
	} else {
		c.setError(err)
	}
}

func (c *BaseControl) setError(err error) {
	if newErr := err.Error(); c.LastError != newErr {
		log.Println("Control err:", newErr)
		c.LastError = newErr
	}
}

func (c *BaseControl) transformValue() interface{} {
	if c.Transform != "" {
		if c.vm == nil {
			c.vm = otto.New()
		}

		c.vm.Set("v", c.value)
		result, err := c.vm.Run(c.Transform)
		newVal, _ := result.Export()

		if err != nil {
			c.setError(err)
		}

		return newVal
	} else {
		return c.value
	}
}

/*
 * Fixed Value Control
 */

type FixedValueControl struct {
	BaseControl
	Value interface{}
}

func (c *FixedValueControl) Apply(effect interface{}) {
	c.BaseControl.value = c.Value
	c.BaseControl.Apply(effect)
}

/*
 * Tween Control
 */

type TweenControl struct {
	BaseControl

	Function string
	Min, Max int
	Speed    float64
}

func (c *TweenControl) Apply(effect interface{}) {
	c.BaseControl.value = round(c.getFunction()(float64(c.Min), float64(c.Max), c.Speed))
	c.BaseControl.Apply(effect)
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

/*
 * Midi Control
 */

type MidiControl struct {
	BaseControl

	Device        int
	Status, Data1 int64

	Min, Max float64
	stop     chan struct{}
}

func (c *MidiControl) Init() {
	c.BaseControl.Init()

	if devices := midi.Devices(); len(devices) > c.Device {
		midiMsgs, stop := midi.StreamMessages(devices[c.Device])
		c.stop = stop
		go c.readValues(midiMsgs)
	} else {
		log.Println("MidiControl", "device not found", "| id:", c.Device)
	}
}

func (c *MidiControl) Destroy() {
	if c.stop != nil {
		c.stop <- struct{}{}
	}
}

func (c *MidiControl) Apply(effect interface{}) {
	c.BaseControl.Apply(effect)
}

func (c *MidiControl) readValues(msgs chan midi.MidiMessage) {
	for msg := range msgs {
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

type OscControl struct {
	BaseControl

	Address  string
	Argument int

	Debug interface{}
	stop  chan struct{}
}

func (c *OscControl) Init() {
	c.BaseControl.Init()

	var stream chan osc.OscMessage
	stream, c.stop = osc.StreamMessages()

	go func() {
		for msg := range stream {
			if msg.Address == c.Address && len(msg.Arguments) > c.Argument {
				c.BaseControl.value = msg.Arguments[c.Argument]
				c.Debug = c.value
			}
		}
	}()
}

func (c *OscControl) Destroy() {
	c.stop <- struct{}{} // signals osc stream to quit
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
	case "BaseControl":
		return &BaseControl{}
	}
	return NullControl{}
}

/*
 * Reflection functions
 */

func getField(effect interface{}, path string) (reflect.Value, error) {
	parts := strings.Split(path, ".")
	field := getFieldPart(reflect.ValueOf(effect).Elem(), parts[0])

	for i := 1; i < len(parts); i++ {
		if field.Kind() == reflect.Struct {
			field = getFieldPart(field, parts[i])
		} else {
			return field, errors.New("field not found: " + path)
		}
	}

	if !field.IsValid() {
		return field, errors.New("field not found: " + path)
	}

	return field, nil
}

func getFieldPart(field reflect.Value, part string) reflect.Value {
	tmp := strings.Split(part, "[")
	name, indexes := tmp[0], tmp[1:]

	field = field.FieldByName(name)

	for j := 0; j < len(indexes); j++ {
		index := strings.Split(indexes[j], "]")[0] // discard trailing ]
		i, _ := strconv.Atoi(index)

		field = field.Index(i)

		// transparently unwrap effect envelopes
		if _, ok := field.Interface().(EffectEnvelope); ok {
			field = reflect.Indirect(field.FieldByName("Effect").Elem())
		}
	}

	return field
}

func setValue(field reflect.Value, newVal interface{}) error {
	if err, v := convertValue(newVal, field.Type()); err == nil {
		switch field.Interface().(type) {
		case int:
			field.SetInt(int64(v.(int)))
		case uint8:
			field.SetUint(uint64(v.(uint8)))
		case float64:
			field.SetFloat(float64(v.(float64)))
		case bool:
			field.SetBool(bool(v.(bool)))
		default:
			return errors.New("can't set field (unknown type: " + field.Type().String())
		}
		return nil
	} else {
		return err
	}
}

func convertValue(val interface{}, t reflect.Type) (error, interface{}) {
	if reflect.TypeOf(val).ConvertibleTo(t) {
		return nil, reflect.ValueOf(val).Convert(t).Interface()
	}
	return errors.New("can't convert " + reflect.TypeOf(val).String() + "->" + t.String()), nil
}
