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

// cycle between min and max values, restarting at min
// once max has been reached e.g 1-2-3-4-1-2-3-4-1-2-3
func CycleBetween(min, max, hertz float64) float64 {
	d := time.Since(timeBegan).Seconds()
	return sawtoothWave(d*hertz)*(max-min) + min
}

//   ^   ^   ^   ^   ^   ^   ^   ^   ^
//  / \ / \ / \ / \ / \ / \ / \ / \ / \
// /   v   v   v   v   v   v   v   v   v
// returns values between 0..1, period=2
func triangleWave(x float64) float64 {
	x = math.Abs(x - 1) // correct the phase (start at 0)
	return math.Abs(math.Mod(x, 2) - 1)
}

//   /|  /|  /|  /|  /|  /|  /|  /|  /|
//  / | / | / | / | / | / | / | / | / |
// /  |/  |/  |/  |/  |/  |/  |/  |/  |/
// returns values between 0..1, period=1
func sawtoothWave(x float64) float64 {
	return math.Mod(x, 1)
}

//      ooo         ooo        ooo
//    o     o     o     o    o     o
//   o        ooo         ooo        ooo
// returns values between 0..1, period=1
func sinusoidWave(x float64) float64 {
	return math.Sin(2*math.Pi*x)/2 + 0.5
}

// simple rounding function
// still not present in stdlib (Go 1.9)
func round(x float64) int {
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}
