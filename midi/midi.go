package midi

import (
	"github.com/rakyll/portmidi"
	"log"
)

func init() {
	portmidi.Initialize()
}

type Device struct {
	Name string
	id   portmidi.DeviceID
}

type MidiMessage struct {
	Timestamp int64
	Status    int64
	Data1     int64
	Data2     int64
}

func Devices() []Device {
	var devices []Device
	count := portmidi.CountDevices()
	for i := 0; i < count; i++ {
		id := portmidi.DeviceID(i)
		info := portmidi.Info(id)
		devices = append(devices, Device{
			Name: info.Name,
			id:   id,
		})
	}
	return devices
}

func StreamMessages(device Device) chan MidiMessage {
	in, err := portmidi.NewInputStream(device.id, 1024)
	if err != nil {
		log.Fatal(err)
	}

	inChan := in.Listen()
	outChan := make(chan MidiMessage)

	go func() {
		for event := range inChan {
			outChan <- MidiMessage{
				Timestamp: int64(event.Timestamp),
				Status:    event.Status,
				Data1:     event.Data1,
				Data2:     event.Data2,
			}
		}
	}()

	return outChan
}
