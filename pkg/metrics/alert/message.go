package alert

import (
	"fmt"
	"time"
)

const (
	HighTraffic = iota
	Resolved
)

type Msg struct {
	Type  int
	Value float64
	When  time.Time
}

func (a *Msg) String() string {
	return fmt.Sprintf(
		"High traffic generated an alert - hits = %.2f, triggered at %s",
		a.Value, a.When.Format(time.RFC3339),
	)
}
