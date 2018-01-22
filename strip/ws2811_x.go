// +build !ws2811

package strip

import (
	"log"
	"time"
)

func (s *LEDStrip) RunWs2811(gpioPin int, refreshInterval time.Duration, maxBrightness int) {
	log.Println("WS2811 unsupported in this build, see wiki for info")
}
