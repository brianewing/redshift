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
	count := portmidi.CountDevices()
	devices := make([]Device, count)
	for i := 0; i < count; i++ {
		id := portmidi.DeviceID(i)
		info := portmidi.Info(id)
		devices[i] = Device{
			Name: info.Name,
			id:   id,
		}
	}
	return devices
}

func StreamMessages(device Device) (msgs chan MidiMessage, done chan struct{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if streams[device.id] == nil {
		streams[device.id] = []chan MidiMessage{}
		openStream(device)
	}

	msgs = make(chan MidiMessage)
	streams[device.id] = append(streams[device.id], msgs)

	done = make(chan struct{})
	go waitForDone(msgs, done, device)

	return
}

func waitForDone(msgs chan MidiMessage, done chan struct{}, device Device) {
	<-done
	mutex.Lock()
	for i, c := range streams[device.id] {
		if c == msgs {
			streams[device.id] = append(streams[device.id][:i], streams[device.id][i+1:]...)
			close(c)
		}
	}
	mutex.Unlock()
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
