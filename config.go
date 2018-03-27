package main

import (
	"github.com/brianewing/redshift/effects"
	"io/ioutil"
	"log"
)

func getEffects() effects.EffectSet {
	effects, _ := loadEffects()
	return effects
}

func loadEffects() (effects.EffectSet, error) {
	if *pathToEffectsJson != "" {
		return loadEffectsJson(*pathToEffectsJson)
	} else if *pathToEffectsYaml != "" {
		return loadEffectsYaml(*pathToEffectsYaml)
	} else {
		return nil, nil
	}
}

func loadEffectsJson(path string) (effects.EffectSet, error) {
	if bytes, err := ioutil.ReadFile(path); err == nil {
		return effects.UnmarshalJSON(bytes)
	} else {
		log.Fatalln("Could not load effects json", err)
		return nil, err
	}
}

func loadEffectsYaml(path string) (effects.EffectSet, error) {
	if bytes, err := ioutil.ReadFile(path); err == nil {
		return effects.UnmarshalYAML(bytes)
	} else {
		log.Fatalln("Could not load effects yaml", err)
		return nil, err
	}
}

func writeEffectsJson(dest string, effects_ effects.EffectSet) error {
	bytes, err := effects.MarshalJSON(effects_)
	if err != nil {
		log.Fatalln("Could not write effects json", "(marshall error)", err)
		return err
	}
	return ioutil.WriteFile(dest, bytes, 0644)
}

func writeEffectsYaml(dest string, effects_ effects.EffectSet) error {
	bytes, err := effects.MarshalYAML(effects_)
	if err != nil {
		log.Fatalln("Could not write effects yaml", "(marshall error)", err)
		return err
	}
	return ioutil.WriteFile(dest, bytes, 0644)
}
