package main

import (
	"bufio"
	"encoding/json"
	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/osc"
	"github.com/robertkrimen/otto"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
)

func RunReplServer(address string, animator *animator.Animator) error {
	s, _ := net.Listen("tcp", address)
	for {
		if client, err := s.Accept(); err == nil {
      go repl(animator, client)
    } else {
      return err
    }
	}
  return nil
}

func repl(a *animator.Animator, client io.ReadWriter) {
	scanner := bufio.NewScanner(client)
	jsVm := otto.New()

	print := func(response string) {
		io.WriteString(client, response)
	}

	println := func(response string) {
		io.WriteString(client, response+"\n")
	}

	for scanner.Scan() {
		input := scanner.Text()
		words := strings.Split(input, " ")

		cmd := words[0]
		tail := input[strings.Index(input, " ")+1:]

		switch cmd {
		case "h", "help", "?":
			println("(e) effects, (e.y) effects.yaml, (e.j) effects.json, (t) types, (a) add, (p) pop, (s) shift, (n) count")
			if len(words) > 1 {
				println("dev: (g|g.n) goroutines, (m|m.n) mutexes, (h|h.n) heap, (osc|osc.c) osc, (.) eval, (q) quit")
			}

		case "gc!":
			runtime.GC()

		case "g", "goroutines":
			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		case "g.n":
			println(strconv.Itoa(pprof.Lookup("goroutine").Count()) + " goroutines")

		case "m", "mutexes":
			pprof.Lookup("mutex").WriteTo(os.Stdout, 1)
		case "m.n":
			println(strconv.Itoa(pprof.Lookup("mutex").Count()) + " mutexes")

		case "heap":
			pprof.Lookup("heap").WriteTo(os.Stdout, 1)
		case "h.n":
			println(strconv.Itoa(pprof.Lookup("heap").Count()) + " heaps")

		case ".", "eval":
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
			json, _ := json.MarshalIndent(a.Effects, "", "  ")
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

		case "u", "unshift":
			var newEffect effects.EffectEnvelope
			if err := newEffect.UnmarshalJSON([]byte(tail)); err != nil {
				println(err.Error())
			} else {
				newEffect.Init()
				a.Effects = append(effects.EffectSet{newEffect}, a.Effects...)
			}

		case "n", "count":
			println(strconv.Itoa(len(a.Effects)) + " effects")

		case "osc":
			summary := osc.Summary()
			oscJson, _ := json.MarshalIndent(summary, "", "  ")
			print("Summary: ")
			println(string(oscJson))

		case "osc.c":
			osc.ClearSummary()

		default:
			println("?")

		case "":
		}
		print("> ")
	}

	println("repl done")
}
