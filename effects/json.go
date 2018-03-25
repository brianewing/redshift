package effects

import (
	json "github.com/brianewing/redshift/customjson"
	"reflect"
)

// effects are json encoded into envelopes

type jsonEnvelope struct {
	Type   string
	Params interface{}
}

type unmarshalFormat struct {
	jsonEnvelope
	Params json.RawMessage
}

func MarshalJSON(effects []Effect) ([]byte, error) {
	envelopes := make([]jsonEnvelope, len(effects))

	for i, effect := range effects {
		envelopes[i] = jsonEnvelope{
			Type:   reflect.TypeOf(effect).Elem().Name(),
			Params: effect,
		}
	}

	return json.Marshal(envelopes)
}

func UnmarshalJSON(s []byte) ([]Effect, error) {
	var envelopes []unmarshalFormat
	json.Unmarshal(s, &envelopes)

	effects := make([]Effect, len(envelopes))
	for i, envelope := range envelopes {
		effects[i] = NewByName(envelope.Type)
		json.Unmarshal(envelope.Params, &effects[i])
	}

	return effects, nil
}

/*
 * Effect Sets (e.g. layer.Effects, mirror.Effects)
 */

func (set *EffectSet) MarshalJSON() ([]byte, error) {
	return MarshalJSON(*set)
}

func (set *EffectSet) UnmarshalJSON(b []byte) error {
	if effects, err := UnmarshalJSON(b); err == nil {
		for _, effect := range effects {
			*set = append(*set, effect)
		}
		return nil
	} else {
		return err
	}
}
