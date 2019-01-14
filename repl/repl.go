package repl

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/brianewing/redshift/animator"
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/osc"
	"github.com/robertkrimen/otto"
)

var serverStartTime = time.Now()

func RunReplServer(address string, animator *animator.Animator) error {
	s, _ := net.Listen("tcp", address)
	for {
		if client, err := s.Accept(); err == nil {
			go Run(animator, client, "> ")
		} else {
			return err
		}
	}
	return nil
}

func Run(a *animator.Animator, client io.ReadWriter, prompt string) {
	scanner := bufio.NewScanner(client)
	jsVm := otto.New()

	print := func(response string) {
		io.WriteString(client, response)
	}

	println := func(response string) {
		io.WriteString(client, response+"\n")
	}

	var stopLoggingFrameRate chan struct{}

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
			pprof.Lookup("goroutine").WriteTo(client, 2)
		case "g.n":
			println(strconv.Itoa(pprof.Lookup("goroutine").Count()) + " goroutines")

		case "m", "mutexes":
			pprof.Lookup("mutex").WriteTo(client, 1)
		case "m.n":
			println(strconv.Itoa(pprof.Lookup("mutex").Count()) + " mutexes")

		case "heap":
			pprof.Lookup("heap").WriteTo(client, 1)
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
			if len(a.Effects) > 0 {
				a.Effects[len(a.Effects)-1].Destroy()
				a.Effects = a.Effects[:len(a.Effects)-1]
			}

		case "s", "shift":
			if len(a.Effects) > 0 {
				a.Effects[0].Destroy()
				a.Effects = a.Effects[1:]
			}

		case "u", "unshift":
			var newEffect effects.EffectEnvelope
			if err := newEffect.UnmarshalJSON([]byte(tail)); err != nil {
				println(err.Error())
			} else {
				newEffect.Init()
				a.Effects = append(effects.EffectSet{newEffect}, a.Effects...)
			}

		case "uptime":
			println(time.Now().Sub(serverStartTime).String())

		case "n", "count":
			println(strconv.Itoa(len(a.Effects)) + " effects")

		case "osc":
			summary := osc.Summary()
			oscJson, _ := json.MarshalIndent(summary, "", "  ")
			print("Summary: ")
			println(string(oscJson))

		case "osc.c":
			osc.ClearSummary()

		case "fps":
			if stopLoggingFrameRate != nil {
				stopLoggingFrameRate <- struct{}{}
				stopLoggingFrameRate = nil
			} else {
				stopLoggingFrameRate = logFrameRate(a, client)
			}

		case "q", "quit", "exit":
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		default:
			println("?")

		case "":
		}
		print(prompt)
	}

	// cleanup
	println("cleanup repl")
	if stopLoggingFrameRate != nil {
		stopLoggingFrameRate <- struct{}{}
	}
}

func logFrameRate(animator *animator.Animator, client io.ReadWriter) (stop chan struct{}) {
	stop = make(chan struct{})

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-stop:
				ticker.Stop()
				return
			case <-ticker.C:
				client.Write(append([]byte(animator.Performance.String()), '\n'))
				animator.Performance.Reset()
			}
		}
	}()

	return stop
}
