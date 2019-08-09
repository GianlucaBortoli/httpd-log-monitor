package avg

import (
	"fmt"
	"time"
)

// Avg implements a metric that can return the per-second average of a counter
// over a time frame of a given size
type Avg struct {
	count      float64
	windowSize time.Duration
}

// New returns an average metric object
func New(t time.Duration) (*Avg, error) {
	if t == 0 {
		return nil, fmt.Errorf("cannot use a time window of size %d", t)
	}
	return &Avg{windowSize: t}, nil
}

// IncrBy increments the counter by i
func (a *Avg) IncrBy(i float64) {
	if i <= 0 {
		return
	}
	a.count += i
}

// GetAvgPerSec returns the per-second average of counter over the windowSize
func (a *Avg) GetAvgPerSec() float64 {
	return a.count / a.windowSize.Seconds()
}

// Reset sets the counter back to zero
func (a *Avg) Reset() {
	a.count = 0
}
