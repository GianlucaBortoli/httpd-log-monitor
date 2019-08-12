package alert

import (
	"fmt"
	"time"
)

const (
	highTraffic = iota
	resolved
)

type msg struct {
	Type  int
	Value float64
	When  time.Time
}

func (a *msg) String() string {
	switch a.Type {
	case highTraffic:
		return fmt.Sprintf(
			"High traffic generated an alert - hits = %.2f, triggered at %s",
			a.Value, a.When.Format(time.RFC3339),
		)
	case resolved:
		return fmt.Sprintf(
			"High traffic alert resolved - hits = %.2f, triggered at %s",
			a.Value, a.When.Format(time.RFC3339),
		)
	default:
		return fmt.Sprintf("unknown alert type %d", a.Type)
	}
}
