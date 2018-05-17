package effects

import (
	"encoding/json"
	"reflect"
)

// Effect envelopes

type _effectEnvelope EffectEnvelope

type effectMarshalFormat struct {
	_effectEnvelope // using a type alias prevents json.Marshal from recursing
	Type string
}

type effectUnmarshalFormat struct {
	effectMarshalFormat
	Effect *json.RawMessage // this can only be unpacked once the Type is known
}

func (e *EffectEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal(effectMarshalFormat{
		Type: reflect.TypeOf(e.Effect).Elem().Name(),
		_effectEnvelope: _effectEnvelope(*e),
	})
}

func (e *EffectEnvelope) UnmarshalJSON(b []byte) error {
	var tmp effectUnmarshalFormat
	if err := json.Unmarshal(b, &tmp); err == nil {
		e.Effect = NewByName(tmp.Type)
		e.Controls = tmp.Controls

		if tmp.Effect != nil {
			return json.Unmarshal(*tmp.Effect, &e.Effect)
		} else {
			return nil
		}
	} else {
		return err
	}
}

// Control envelopes

type _controlEnvelope ControlEnvelope

type controlMarshalFormat struct {
	_controlEnvelope // using a type alias prevents json.Marshal from recursing
	Type string
}

type controlUnmarshalFormat struct {
	controlMarshalFormat
	Control *json.RawMessage // this can only be unpacked once the Type is known
}

func (e *ControlEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal(controlMarshalFormat{
		_controlEnvelope: _controlEnvelope(*e),
		Type: reflect.TypeOf(e.Control).Elem().Name(),
	})
}

func (e *ControlEnvelope) UnmarshalJSON(b []byte) error {
	var tmp controlUnmarshalFormat
	if err := json.Unmarshal(b, &tmp); err == nil {
		e.Control = ControlByName(tmp.Type)
		e.Controls = tmp.Controls
		if tmp.Control != nil {
			return json.Unmarshal(*tmp.Control, &e.Control)
		} else {
			return nil
		}
	} else {
		return err
	}
}

// Helper functions

func MarshalJSON(effects EffectSet) ([]byte, error) {
	return json.Marshal(effects)
}

func UnmarshalJSON(b []byte) (EffectSet, error) {
	var effects EffectSet
	if err := json.Unmarshal(b, &effects); err == nil {
		return effects, nil
	} else {
		return nil, err
	}
}
