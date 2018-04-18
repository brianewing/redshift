package midi

import (
	"github.com/rakyll/portmidi"
	"sync"
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
	mutex.Lock()
	defer mutex.Unlock()

	if streams[device.id] == nil {
		streams[device.id] = []chan MidiMessage{}
		openStream(device)
	}

	outChan := make(chan MidiMessage)
	streams[device.id] = append(streams[device.id], outChan)

	return outChan
}

var streams = map[portmidi.DeviceID][]chan MidiMessage{}
var mutex sync.Mutex

func openStream(device Device) error {
	in, err := portmidi.NewInputStream(device.id, 1024)
	if err != nil {
		return err
	}
	inChan := in.Listen()

	go func() {
		for event := range inChan {
			mutex.Lock()

			for _, outChan := range streams[device.id] {
				outChan <- MidiMessage{
					Timestamp: int64(event.Timestamp),
					Status:    event.Status,
					Data1:     event.Data1,
					Data2:     event.Data2,
				}
			}

			mutex.Unlock()
		}
	}()

	return nil
}
