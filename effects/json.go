package effects

import (
	json "github.com/brianewing/redshift/customjson"
	"reflect"
)

// TODO: document this code better!

type myEffectEnvelope EffectEnvelope

type marshalFormat struct {
	myEffectEnvelope // using a type alias prevents json.Marshal from recursing
	Type string
}

type unmarshalFormat struct {
	marshalFormat
	Effect *json.RawMessage // this can only be unpacked once the Type is known
	Controls *json.RawMessage
}

func (e *EffectEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalFormat{
		myEffectEnvelope: myEffectEnvelope(*e),
		Type: reflect.TypeOf(e.Effect).Elem().Name(),
	})
}

func (e *EffectEnvelope) UnmarshalJSON(b []byte) error {
	var tmp unmarshalFormat
	if err := json.Unmarshal(b, &tmp); err == nil {
		e.Effect = NewByName(tmp.Type)
		if tmp.Effect != nil {
			return json.Unmarshal(*tmp.Effect, &e.Effect)
		} else {
			return nil
		}
	} else {
		return err
	}
}

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
