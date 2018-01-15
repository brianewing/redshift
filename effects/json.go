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
 * Layer Effect
 */

type jsonFormatLayer struct {
	Effects json.RawMessage // encoded with MarshalJson()
}

func (e *Layer) MarshalJSON() ([]byte, error) {
	effectsJson, _ := MarshalJSON(e.Effects)
	return json.Marshal(&jsonFormatLayer{Effects: effectsJson})
}

func (e *Layer) UnmarshalJSON(b []byte) error {
	tmp := jsonFormatLayer{}
	if err := json.Unmarshal(b, &tmp); err == nil {
		if effects, err := UnmarshalJSON(tmp.Effects); err == nil {
			e.Effects = effects
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
