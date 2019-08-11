package rate

import (
	"fmt"
	"time"
)

// Rate implements a metric that can return the per-second average of a counter
// over a time frame of a given size
type Rate struct {
	count      float64
	windowSize time.Duration
}

// New returns an average metric object
func New(t time.Duration) (*Rate, error) {
	if t == 0 {
		return nil, fmt.Errorf("cannot have time window of width 0")
	}
	return &Rate{windowSize: t}, nil
}

// IncrBy increments the counter by i
func (r *Rate) IncrBy(i float64) error {
	if i < 0 {
		return fmt.Errorf("cannot increment by negavite number")
	}
	r.count += i
	return nil
}

// GetAvgPerSec returns the per-second average of counter over the windowSize
func (r *Rate) GetAvgPerSec() float64 {
	return r.count / r.windowSize.Seconds()
}

// Reset sets the counter back to zero
func (r *Rate) Reset() {
	r.count = 0
}

func (r *Rate) GetWindowSize() time.Duration {
	return r.windowSize
}

func (r *Rate) GetCount() float64 {
	return r.count
}
