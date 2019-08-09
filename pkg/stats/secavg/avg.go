package secavg

import (
	"fmt"
	"time"
)

// SecAvg implements a metric that can return the per-second average of a counter
// over a time frame of a given size
type SecAvg struct {
	count      float64
	windowSize time.Duration
}

// New returns an average metric object
func New(t time.Duration) (*SecAvg, error) {
	if t == 0 {
		return nil, fmt.Errorf("cannot use a time window of size %d", t)
	}
	return &SecAvg{windowSize: t}, nil
}

// IncrBy increments the counter by i
func (s *SecAvg) IncrBy(i float64) error {
	if i < 0 {
		return fmt.Errorf("cannot increment by %f", i)
	}
	s.count += i
	return nil
}

// GetAvgPerSec returns the per-second average of counter over the windowSize
func (s *SecAvg) GetAvgPerSec() float64 {
	return s.count / s.windowSize.Seconds()
}

// Reset sets the counter back to zero
func (s *SecAvg) Reset() {
	s.count = 0
}

func (s *SecAvg) GetWindowSize() time.Duration {
	return s.windowSize
}
