package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brianewing/redshift/opc"
	"github.com/brianewing/redshift/strip"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var addr = flag.String("addr", "localhost:7890", "opc address")

var clear = flag.Bool("clear", true, "clear output before each frame")
var _return = flag.Bool("return", false, "line feed after each frame")

var leds = flag.Int("leds", 60, "number of leds for virtual streams (use with -yaml)")
var fps = flag.Int("fps", 60, "fps")

var letters = flag.Bool("letters", false, "show letters in output")

var jsonFile = flag.String("json", "", "path to json effects definition")
var yamlFile = flag.String("yaml", "", "path to yaml effects definition")

type OpcConn struct {
	io.ReadWriter
}

func (c OpcConn) ReadMsg() (opc.Message, error) {
	return opc.ReadMessage(c)
}

func (c OpcConn) SendMsg(msg opc.Message) error {
	_, err := c.Write(msg.Bytes())
	return err
}

func (c OpcConn) SendSysEx(channel uint8, command opc.SystemExclusiveCmd, data []byte) error {
	return c.SendMsg(opc.Message{
		Channel: channel,
		Command: 255,
		SystemExclusive: opc.SystemExclusive{
			Command: command,
			Data:    data,
		},
	})
}

func (c OpcConn) OpenStream(channel uint8, desc string) error {
	return c.SendSysEx(channel, opc.CmdOpenStream, []byte(desc))
}

func (c OpcConn) SetStreamFps(channel uint8, fps uint8) error {
	return c.SendSysEx(channel, opc.CmdSetStreamFps, []byte{fps})
}

type OpcClient struct {
	OpcConn
	listeners []chan opc.Message
}

func (c *OpcClient) SubscribeToMsgs() (chan opc.Message, chan struct{}) {
	listener := make(chan opc.Message)
	done := make(chan struct{})
	c.listeners = append(c.listeners, listener)
	go c.removeListenerWhenDone(listener, done)
	return listener, done
}

func (c *OpcClient) WaitForMsg(matchFn func(opc.Message) bool) opc.Message {
	msgs, done := c.SubscribeToMsgs()
	defer func() { done <- struct{}{} }()
	for msg := range msgs {
		if matchFn(msg) {
			return msg
		}
	}
	return opc.Message{}
}

func (c *OpcClient) removeListenerWhenDone(listener chan opc.Message, done chan struct{}) {
	<-done
	for i, l := range c.listeners {
		if l == listener {
			c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
			close(l)
			break
		}
	}
}

func (c *OpcClient) ReadMsgs() {
	for {
		msg, err := c.OpcConn.ReadMsg()
		if err != nil {
			fmt.Println("read err", err)
			return
		}
		for _, l := range c.listeners {
			l <- msg
		}
	}
}

func (c *OpcClient) Effects() []byte {
	go c.SendSysEx(0, opc.CmdSetEffectsStreamFps, []byte{0})
	msg := c.WaitForMsg(func(msg opc.Message) bool {
		return msg.SystemExclusive.Command == opc.CmdSetEffectsJson
	})
	return msg.SystemExclusive.Data
}

func hideCursor() { print("\u001B[?25l") }
func showCursor() { print("\u001B[?25h") }
func clearLine()  { print("\r") }

func clearPreviousLine(chars int) {
	print("\033[F")
	for i := 0; i < chars; i++ {
		print(" ")
	}
	print("\r")
}

func printPixels(buffer strip.Buffer) {
	if *clear {
		clearLine()
	}
	for _, pixel := range buffer {
		r := pixel[0]
		g := pixel[1]
		b := pixel[2]

		fmt.Printf("\033[48;2;%d;%d;%dm", r, g, b)
		fmt.Printf("\033[38;2;%d;%d;%dm", 255-r, 255-g, 255-b) // set foreground color
		fmt.Printf("%s\033[0m", letter())
	}
	if *_return {
		print("\n")
	}
}

func readFile(name string) []byte {
	bytes, _ := ioutil.ReadFile(name)
	return bytes
}

var s = "abcdefghijklmnopqrstuvwxyz0123456789!@£$%^&*()_+}{|\""

// var s = readFile("main.go")
var i int

func letter() string {
	if *letters {
		x := s[i%len(s)]
		i += 1
		return string(x)
	} else {
		return " "
	}
}

func main() {
	flag.Parse()

	if conn, err := net.Dial("tcp", *addr); err != nil {
		fmt.Println("tcp error", err)
	} else {
		opcConn := OpcConn{ReadWriter: conn}

		if *yamlFile == "" && *jsonFile == "" {
			opcConn.OpenStream(0, "strip")
		} else {
			opcConn.OpenStream(0, "virtual "+strconv.Itoa(*leds))
		}

		if *yamlFile != "" {
			opcConn.SendSysEx(0, opc.CmdSetEffectsYaml, readFile(*yamlFile))
		} else if *jsonFile != "" {
			opcConn.SendSysEx(0, opc.CmdSetEffectsJson, readFile(*jsonFile))
		}

		opcConn.SetStreamFps(0, uint8(*fps))

		hideCursor()

		var buffer strip.Buffer

		var opcClient = OpcClient{OpcConn: opcConn}
		var msgs, _ = opcClient.SubscribeToMsgs()

		var running = true

		go opcClient.ReadMsgs()

		var replInterrupt = make(chan bool)

		go func() {
			scanner := bufio.NewScanner(os.Stdin)

			resume := func() {
				if !running {
					clearPreviousLine(len(buffer))
					replInterrupt <- false
				}
			}

			for scanner.Scan() {
				input := scanner.Text()

				words := strings.Split(input, " ")
				cmd := words[0]
				tail := input[strings.Index(input, " ")+1:]

				switch cmd {
				case "h", "help", "?":
					println("(o) open, (e) effects, (r) return, (c) clear-line, (y) yaml, (f) fps, (a) add, (r) resume")
				case "o", "open":
					opcConn.SendSysEx(0, opc.CmdCloseStream, []byte{})
					opcConn.SendSysEx(0, opc.CmdOpenStream, []byte(tail))
					resume()
					buffer = nil
				case "e", "effects":
					println(string(opcClient.Effects()))
				case "a", "add":
					opcConn.SendSysEx(0, opc.CmdAppendEffectsJson, []byte(tail))
				case "p", "pop":
					var effects []interface{}
					json.Unmarshal(opcClient.Effects(), &effects)
					effects = effects[:len(effects)-1]
					jsonBytes, _ := json.Marshal(effects)
					opcConn.SendSysEx(0, opc.CmdSetEffectsJson, jsonBytes)
				case "r", "return":
					*_return = !*_return
					resume()
				case "c", "clear-line":
					*clear = !*clear
					resume()
				case "j", "json":
					definition := readFile(words[1])
					opcConn.SendSysEx(0, opc.CmdSetEffectsJson, definition)
					resume()
				case "y", "yaml":
					definition := readFile(words[1])
					opcConn.SendSysEx(0, opc.CmdSetEffectsYaml, definition)
					resume()
				case "f", "fps":
					newFps, _ := strconv.Atoi(words[1])
					*fps = newFps
					resume()
					opcConn.SetStreamFps(0, 0) // request new frame immediately
				case "", "resume":
					if running {
						replInterrupt <- true
						clearPreviousLine(len(buffer))
					} else {
						resume()
					}
				default:
					println("?")
				}

				print("termclient> ")
			}
		}()

		go restoreCursorOnCtrlC()

		for {
			select {
			case stop := <-replInterrupt:
				running = !stop
				if stop {
					opcConn.SetStreamFps(0, 0)
					showCursor()
				} else {
					opcConn.SetStreamFps(0, uint8(*fps))
					hideCursor()
				}
			case msg := <-msgs:
				if !running {
					continue
				}
				if msg.Command == 0 && err == nil {
					if buffer == nil {
						buffer = strip.NewBuffer(len(msg.Data) / 3)
					}
					msg.WritePixels(buffer)
					printPixels(buffer)
					if *fps == 0 {
						println()
						return
					}
				}
			}
		}
	}

	time.Sleep(1 * time.Second)
	main() // reconnect
}

func restoreCursorOnCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	showCursor()
	println()
	os.Exit(0)
}
