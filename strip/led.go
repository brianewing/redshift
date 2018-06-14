package strip

import (
	"bytes"
	"strconv"
)

type LED []uint8

func (led LED) IsOff() bool {
	return led[0] == 0 && led[1] == 0 && led[2] == 0
}

func (led LED) MarshalJSON() ([]byte, error) {
	var tmp bytes.Buffer
	tmp.WriteRune('[')
	for j, val := range led {
		if j != 0 {
			tmp.WriteRune(',')
		}
		tmp.WriteString(strconv.Itoa(int(val)))
	}
	tmp.WriteRune(']')
	return tmp.Bytes(), nil
}
