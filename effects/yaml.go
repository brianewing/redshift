package effects

import "github.com/ghodss/yaml"

func UnmarshalYAML(s []byte) ([]Effect, error) {
	if json, err := yaml.YAMLToJSON(s); err == nil {
		return UnmarshalJSON(json)
	} else {
		return nil, err
	}
}

func MarshalYAML(effects []Effect) ([]byte, error) {
	if json, err := MarshalJSON(effects); err == nil {
		return yaml.JSONToYAML(json)
	} else {
		return nil, err
	}
}
