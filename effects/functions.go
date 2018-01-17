package effects

import (
	"math"
	"time"
)

var timeBegan = time.Now()

// simple time-based linear oscillator
func OscillateBetween(min, max, hertz float64) float64 {
	d := time.Since(timeBegan).Seconds()
	return triangleWave(d*hertz*2)*(max-min) + min
}

//   ^   ^   ^   ^   ^   ^   ^   ^   ^   ^   ^   ^   ^
//  / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \
// v   v   v   v   v   v   v   v   v   v   v   v   v   v
// returns values between 0..1, period=2 with respect to x
func triangleWave(x float64) float64 {
	return math.Abs(math.Mod(x, 2) - 1)
}

// simple rounding function
// still not present in stdlib (Go 1.9)
func round(x float64) int {
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}
