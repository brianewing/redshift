package main

import (
	"flag"
	"fmt"
	"github.com/brianewing/redshift/server"
	"github.com/brianewing/redshift/strip"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

var addr = flag.String("addr", "localhost:7890", "opc address")

var clear = flag.Bool("clear", true, "clear output before each frame")
var _return = flag.Bool("return", false, "carriage return after each frame")

var leds = flag.Int("leds", 60, "number of leds for virtual streams (use with -yaml)")
var fps = flag.Int("fps", 60, "fps")

var letters = flag.Bool("letters", false, "show letters in output")

var json = flag.String("json", "", "path to json effects definition")
var yaml = flag.String("yaml", "", "path to yaml effects definition")

type OpcConn struct {
	io.ReadWriter
}

func (c OpcConn) ReadMsg() (server.OpcMessage, error) {
	return server.ReadOpcMessage(c)
}

func (c OpcConn) SendMsg(msg server.OpcMessage) error {
	_, err := c.Write(msg.Bytes())
	return err
}

func (c OpcConn) SendSysEx(channel uint8, command server.SystemExclusiveCmd, data []byte) error {
	return c.SendMsg(server.OpcMessage{
		Channel: channel,
		Command: 255,
		SystemExclusive: server.SystemExclusive{
			Command: command,
			Data:    data,
		},
	})
}

func (c OpcConn) OpenStream(channel uint8, desc string) error {
	return c.SendSysEx(channel, server.CmdOpenStream, []byte(desc))
}

func (c OpcConn) SetStreamFps(channel uint8, fps uint8) error {
	return c.SendSysEx(channel, server.CmdSetStreamFps, []byte{fps})
}

func hideCursor() { print("\u001B[?25l") }
func clearLine()  { print("\r") }

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

		if *yaml == "" && *json == "" {
			opcConn.OpenStream(0, "strip")
		} else {
			opcConn.OpenStream(0, "virtual "+strconv.Itoa(*leds))
		}

		if *yaml != "" {
			opcConn.SendSysEx(0, server.CmdSetEffectsYaml, readFile(*yaml))
		} else if *json != "" {
			opcConn.SendSysEx(0, server.CmdSetEffectsJson, readFile(*json))
		}

		opcConn.SetStreamFps(0, uint8(*fps))

		hideCursor()

		var buffer strip.Buffer

		for {
			msg, err := opcConn.ReadMsg()
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
			} else {
				fmt.Println("read error", err)
				break
			}
		}
	}

	time.Sleep(1 * time.Second)
	main() // reconnect
}
