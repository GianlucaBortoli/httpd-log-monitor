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
	switch a.Type {
	case HighTraffic:
		return fmt.Sprintf(
			"High traffic generated an alert - hits = %.2f, triggered at %s",
			a.Value, a.When.Format(time.RFC3339),
		)
	case Resolved:
		return fmt.Sprintf(
			"High traffic alert resolved - hits = %.2f, triggered at %s",
			a.Value, a.When.Format(time.RFC3339),
		)
	default:
		return fmt.Sprintf("unknown alert type %d", a.Type)
	}
}
