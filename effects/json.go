package effects

import (
	"encoding/json"
	"reflect"
)

// todo: document how this works!

type jsonEnvelope struct {
	Type string
	Params interface{}
}

type unmarshalFormat struct {
	jsonEnvelope
	Params json.RawMessage
}

func MarshalJson(effects []Effect) ([]byte, error) {
	envelopes := make([]jsonEnvelope, len(effects))

	for i, effect := range effects {
		envelopes[i] = jsonEnvelope{
			Type: reflect.TypeOf(effect).Elem().Name(),
			Params: effect,
		}
	}

	return json.Marshal(envelopes)
}

func UnmarshalJson(s []byte) ([]Effect, error) {
	var envelopes []unmarshalFormat
	json.Unmarshal(s, &envelopes)

	effects := make([]Effect, len(envelopes))
	for i, envelope := range envelopes {
		effects[i] = newEffectByName(envelope.Type)
		json.Unmarshal(envelope.Params, &effects[i])
	}

	return effects, nil
}

/*
 * Combine Effect
 */

type jsonFormatCombine struct {
	Effects json.RawMessage // encoded with MarshalJson()
}

func (e *Combine) MarshalJSON() ([]byte, error) {
	effectsJson, _ := MarshalJson(e.Effects)
	return json.Marshal(&jsonFormatCombine{Effects: effectsJson})
}

func (e *Combine) UnmarshalJSON(b []byte) error {
	tmp := jsonFormatCombine{}
	if err := json.Unmarshal(b, &tmp); err == nil {
		if effects, err := UnmarshalJson(tmp.Effects); err == nil {
			e.Effects = effects
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
