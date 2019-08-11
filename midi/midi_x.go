// +build !midi

package midi

type Device struct {
	Name string
}

type MidiMessage struct {
	Timestamp int64
	Status    int64
	Data1     int64
	Data2     int64
}

func Devices() []Device {
	return []Device{}
}

func StreamMessages(device Device) (msgs chan MidiMessage, done chan struct{}) {
	return nil, nil
}
