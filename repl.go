package main

import (
	"bufio"
	"encoding/json"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/robertkrimen/otto"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func repl(a *animator.Animator) {
	scanner := bufio.NewScanner(os.Stdin)
	jsVm := otto.New()

	for scanner.Scan() {
		input := scanner.Text()
		words := strings.Split(input, " ")

		cmd := words[0]
		tail := input[strings.Index(input, " ")+1:]

		switch cmd {
		case "h", "help", "?":
			println("(e) effects, (e.y) effects.yaml, (e.j) effects.json, (t) types, (a) add, (p) pop, (s) shift, (n) count")

		case ".":
			jsVm.Set("a", a)
			jsVm.Set("fx", a.Effects)

			value, _ := jsVm.Run(tail)
			result, _ := value.Export()

			jsonString, _ := json.Marshal(result)
			println(string(jsonString))

		case "e", "effects":
			types := make([]string, len(a.Effects))
			for i, e := range a.Effects {
				types[i] = reflect.TypeOf(e.Effect).Elem().Name()
			}
			println(strings.Join(types, ", "))

		case "e.y", "effects.yaml":
			yaml, _ := effects.MarshalYAML(a.Effects)
			println(string(yaml))

		case "e.j", "effects.json":
			json, _ := effects.MarshalJSON(a.Effects)
			println(string(json))

		case "t", "types":
			println(strings.Join(effects.Names(), ", "))

		case "a", "add":
			var newEffect effects.EffectEnvelope
			if err := newEffect.UnmarshalJSON([]byte(tail)); err != nil {
				println(err.Error())
			} else {
				newEffect.Init()
				a.Effects = append(a.Effects, newEffect)
			}

		case "p", "pop":
			a.Effects[len(a.Effects)-1].Destroy()
			a.Effects = a.Effects[:len(a.Effects)-1]

		case "s", "shift":
			a.Effects[0].Destroy()
			a.Effects = a.Effects[1:]

		case "n", "count":
			println(strconv.Itoa(len(a.Effects)) + " effects")

		default:
			println("?")

		case "":
		}
		print("> ")
	}
}
