package effects

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/brianewing/redshift/midi"
	"github.com/brianewing/redshift/osc"
	"github.com/robertkrimen/otto"
)

type Control interface {
	Apply(effect interface{})
}

type ControlEnvelope struct {
	Control
	Controls ControlSet // recursively controllable
	Disabled bool
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
	if e.Disabled {
		return
	}
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
	transform *otto.Script
	t2        *otto.Object
}

func (c *BaseControl) Init() {
	c.value = c.Initial
	c.LastError = ""
	c.initTransform()
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

func (c *BaseControl) initTransform() {
	if c.Transform != "" {
		c.vm = otto.New()

		var err error
		c.transform, err = c.vm.Compile("", c.Transform)

		if err != nil {
			c.setError(err)
		}
	}
}

func (c *BaseControl) transformValue() interface{} {
	if c.transform != nil {
		if c.vm == nil {
			c.vm = otto.New()
		}

		c.vm.Set("v", c.value)
		result, err := c.vm.Run(c.transform)

		if err != nil {
			c.setError(err)
			return c.value
		}

		newVal, _ := result.Export()
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
	case "tri", "triangle", "linear":
		return OscillateBetween
	case "saw", "sawtooth":
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
 * Time Control
 */

type TimeControl struct {
	BaseControl

	Format string // e.g. day (0-6), hour (0-23), minute (0-59), second (0-59), unix
}

func (c *TimeControl) Apply(effect interface{}) {
	c.BaseControl.value = c.getValue()
	c.BaseControl.Apply(effect)
}

func (c *TimeControl) getValue() interface{} {
	switch c.Format {
	case "unix":
		return float64(time.Now().UnixNano()) / float64(time.Second)
	case "day", "days":
		return int(time.Now().Weekday())
	case "hour", "hours":
		return int(time.Now().Hour())
	case "minute", "minutes":
		return int(time.Now().Minute())
	case "second", "seconds":
	}
	return int(time.Now().Second())
}

/*
 * Audio Control
 */

type AudioControl struct{}

func (c AudioControl) Apply(interface{}) {}

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
		return &TweenControl{Max: 255}
	case "MidiControl":
		return &MidiControl{}
	case "AudioControl":
		return &AudioControl{}
	case "TimeControl":
		return &TimeControl{}
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

// getField looks up string paths like "Effects[0].Color[0]" on a given struct using reflection
// and returns a reflect.Value for the target field, or an error
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

// getFieldPart looks up a path segment (e.g. 'Name', or 'Effects[0]') under a reflect.Value
// Essentially it wraps reflect.Value.FieldByName() with support for indicies
//
// Additionally, if it encounters an EffectEnvelope while deindexing, it will unwrap the Effect
// inside and return it, so that users don't have to know about the envelope wrapper
func getFieldPart(field reflect.Value, part string) reflect.Value {
	tmp := strings.Split(part, "[")
	name, indexes := tmp[0], tmp[1:]

	field = field.FieldByName(name)

	// dereference pointers, e.g. *Blend => Blend
	// if field.Kind() == reflect.Ptr {
	// field = field.Elem()
	// }

	// follow any indicies
	for j := 0; j < len(indexes); j++ {
		if field.Kind() != reflect.Slice {
			continue
		}

		tmp := strings.Split(indexes[j], "]")[0] // discard trailing ]
		i, _ := strconv.Atoi(tmp)

		field = field.Index(i)

		// if the new field is an EffectEnvelope, unwrap the Effect so that users
		// can reference sub-effect fields without having to know about the intermediate layer
		// (i.e. `Effects[0].SomeParam` instead of `Effects[0].Effect.SomeParam`)
		if _, ok := field.Interface().(EffectEnvelope); ok {
			field = field.FieldByName("Effect").Elem()
		}
	}

	return reflect.Indirect(field)
}

func setValue(field reflect.Value, newVal interface{}) error {
	if v, err := convertValue(newVal, field.Type()); err == nil {
		switch field.Interface().(type) {
		case int:
			field.SetInt(int64(v.(int)))
		case uint8:
			field.SetUint(uint64(v.(uint8)))
		case uint16:
			field.SetUint(uint64(v.(uint16)))
		case float64:
			field.SetFloat(float64(v.(float64)))
		case bool:
			field.SetBool(bool(v.(bool)))
		case string:
			field.SetString(string(v.(string)))
		default:
			return errors.New("can't set field (unknown type: " + field.Type().String() + ")")
		}
		return nil
	} else {
		return err
	}
}

func convertValue(val interface{}, newType reflect.Type) (interface{}, error) {
	t := reflect.TypeOf(val)
	if t == nil {
		return nil, errors.New("value is nil")
	} else if t.ConvertibleTo(newType) {
		return reflect.ValueOf(val).Convert(newType).Interface(), nil
	}
	return nil, errors.New("can't convert " + t.String() + "->" + newType.String())
}
