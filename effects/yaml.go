package effects

import "github.com/ghodss/yaml"

func UnmarshalYAML(s []byte) (EffectSet, error) {
	if json, err := yaml.YAMLToJSON(s); err == nil {
		return UnmarshalJSON(json)
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
