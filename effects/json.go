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
	Params interface{}
}

type unmarshalFormat struct {
	marshalFormat
	Params *json.RawMessage // this can only be unpacked once the Type is known
}

func (e *EffectEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalFormat{
		Type: reflect.TypeOf(e.Effect).Elem().Name(),
		Params: e.Effect,
	})
}

func (e *EffectEnvelope) UnmarshalJSON(b []byte) error {
	var tmp unmarshalFormat
	if err := json.Unmarshal(b, &tmp); err == nil {
		e.Effect = NewByName(tmp.Type)
		if tmp.Params != nil {
			return json.Unmarshal(*tmp.Params, &e.Effect)
		} else {
			return nil
		}
	} else {
		return err
	}
}

func MarshalJSON(effects []Effect) ([]byte, error) {
	envelopes := make([]EffectEnvelope, len(effects))
	for i, effect := range effects {
		envelopes[i].Effect = effect
	}
	return json.Marshal(envelopes)
}

func UnmarshalJSON(b []byte) ([]Effect, error) {
	var envelopes []EffectEnvelope
	if err := json.Unmarshal(b, &envelopes); err == nil {
		effects := make([]Effect, len(envelopes))
		for i, envelope := range envelopes {
			effects[i] = envelope.Effect
		}
		return effects, err
	} else {
		return nil, err
	}
}

/*
 * Effect Sets (e.g. layer.Effects, mirror.Effects)
 */

func (set *EffectSet) MarshalJSON() ([]byte, error) {
	envelopes := make([]EffectEnvelope, len(*set))
	for i, effect := range *set {
		envelopes[i] = EffectEnvelope{Effect: effect}
	}
	return json.Marshal(envelopes)
}

func (set *EffectSet) UnmarshalJSON(b []byte) error {
	var envelopes []EffectEnvelope
	if err := json.Unmarshal(b, &envelopes); err == nil {
		for _, envelope := range envelopes {
			*set = append(*set, envelope.Effect)
		}
		return nil
	} else {
		return err
	}
}
