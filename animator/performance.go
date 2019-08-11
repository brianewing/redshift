package animator

import (
	"fmt"
	"sync"
	"time"
)

type Performance struct {
	frameCount int
	durations  []time.Duration

	startTime time.Time
	mutex     sync.Mutex
}

func NewPerformance() *Performance {
	p := &Performance{}
	p.Reset()
	return p
}

func (p *Performance) Reset() {
	p.mutex.Lock()
	p.frameCount = 0
	p.durations = []time.Duration{}
	p.startTime = time.Now()
	p.mutex.Unlock()
}

func (p *Performance) Tick(t time.Duration) {
	p.mutex.Lock()
	p.frameCount += 1
	p.durations = append(p.durations, t)
	p.mutex.Unlock()
}

func (p *Performance) String() string {
	min, max, avg := p.MinMaxAvg()
	fps := p.FPS()
	return fmt.Sprintf("fps %.2f | %v | %v | %v", fps, min, max, avg)
}

func (p *Performance) FPS() (fps float64) {
	p.mutex.Lock()
	fps = float64(p.frameCount) / float64(time.Now().Sub(p.startTime)) * float64(time.Second)
	p.mutex.Unlock()
	return
}

func (p *Performance) MinMaxAvg() (time.Duration, time.Duration, time.Duration) {
	return min(p.durations), max(p.durations), avg(p.durations)
}

// math functions

func min(durations []time.Duration) (m time.Duration) {
	if len(durations) > 0 {
		m = durations[0]
		for _, d := range durations {
			if d < m {
				m = d
			}
		}
	}
	return
}

func max(durations []time.Duration) (m time.Duration) {
	if len(durations) > 0 {
		m = durations[0]
		for _, d := range durations {
			if d > m {
				m = d
			}
		}
	}
	return
}

func avg(durations []time.Duration) time.Duration {
	if len(durations) > 0 {
		var sum time.Duration
		for i := 0; i < len(durations); i++ {
			sum += durations[i]
		}
		return sum / time.Duration(len(durations))
	} else {
		return 0
	}
}
