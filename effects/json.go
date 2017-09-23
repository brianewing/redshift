package effects

import (
	"encoding/json"
	"reflect"
	"github.com/luci/go-render/render"
)

type jsonFormat struct {
	Type string
	Params interface{}
}

func jsonEnvelope(effect Effect) jsonFormat {
	t := reflect.TypeOf(effect).Elem().Name()
	return jsonFormat{Type: t, Params: effect}
}

func MarshalJson(effects []Effect) ([]byte, error) {
	envelopes := make([]jsonFormat, len(effects))

	for i, effect := range effects {
		envelopes[i] = jsonEnvelope(effect)
	}

	return json.Marshal(envelopes)
}

type unmarshalFormat struct {
	jsonFormat
	Params json.RawMessage
}

func UnmarshalJson(s []byte) ([]Effect, error) {
	var envelopes []unmarshalFormat
	json.Unmarshal(s, &envelopes)

	println(render.Render(envelopes))

	return []Effect{}, nil
}
