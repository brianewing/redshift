package server

import (
	"github.com/kellydunn/go-opc"
	"redshift/strip"
	"time"
)

type Device struct {
	Strip *strip.LEDStrip
	channel uint8
	onUpdate chan int
}

func (d *Device) Write(m *opc.Message) error {
	bytes := m.ByteArray()

	if bytes[1] != 0 {
		return nil
	}

	d.Strip.Sync.Lock()
	for i, val := range bytes[4:] {
		d.Strip.Buffer[i / 3][i % 3] = int(val)
	}
	d.onUpdate <- 1
	d.Strip.Sync.Unlock()

	return nil
}

func (d *Device) Channel() uint8 {
	return 1
}

func RunOpcServer(strip *strip.LEDStrip, onUpdate chan int) {
	s := opc.NewServer()
	s.RegisterDevice(&Device{Strip: strip, onUpdate: onUpdate})

	go s.ListenOnPort("tcp", "localhost:7890")
	go s.Process()

	time.Sleep(1 * time.Second)

	//c := opc.NewClient()
	//c.Connect("tcp", "localhost:7890")
	//m := opc.NewMessage(0)
	//
	//m.SetPixelColor(1, 250, 251, 252)
	//m.SetLength(6)
	//fmt.Println("Valid", m.IsValid())
	//
	//// Send the message!
	//c.Send(m)
}
