package effects

import "github.com/ghodss/yaml"

func UnmarshalYAML(s []byte) (EffectSet, error) {
	if bytes, err := yaml.YAMLToJSON(s); err == nil {
		return UnmarshalJSON(bytes)
	} else {
		return nil, err
	}
}

func MarshalYAML(effects EffectSet) ([]byte, error) {
	if json, err := MarshalJSON(effects); err == nil {
		return yaml.JSONToYAML(json)
	} else {
		return nil, err
	}
}
