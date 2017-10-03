package strip

func (s *LEDStrip) SerializeBytes() []byte {
	return SerializeBufferBytes(s.Buffer)
}

func (s *LEDStrip) UnserializeBytes(bytes []byte) {
	UnserializeBufferBytes(s.Buffer, bytes)
}

func SerializeBufferBytes(buffer [][]uint8) []byte {
	bytes := make([]byte, len(buffer) * 3)
	for i, led := range buffer {
		y := i * 3
		bytes[y] = led[0]
		bytes[y+1] = led[1]
		bytes[y+2] = led[2]
	}
	return bytes
}

func UnserializeBufferBytes(buffer [][]uint8, bytes []byte) {
	for i, val := range bytes {
		if len(buffer) == i / 3 {
			break
		}
		buffer[i / 3][i % 3] = val
	}
}