package strip

import (
	"log"
	"time"
	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

func (s *LEDStrip) RunWs2811(gpioPin int, refreshInterval time.Duration, maxBrightness int) {
	if err := ws2811.Init(gpioPin, s.Size, maxBrightness); err != nil {
		log.Println("Ws2811", "init error", err)
		return
	}

	for {
		s.Lock()
		for i, color := range s.Buffer {
			color32 := uint32((uint32(0) << 24) | (uint32(color[0]) << 16) | (uint32(color[1]) << 8) | uint32(color[2]))
			ws2811.SetLed(i, color32)
		}
		s.Unlock()
		ws2811.Render()
		time.Sleep(refreshInterval)
	}
}
