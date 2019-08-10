package alert

import (
	"fmt"
	"time"
)

type alertMsg struct {
	Value float64
	When  time.Time
}

func (a *alertMsg) String() string {
	return fmt.Sprintf(
		"High traffic generated an alert - hits = %.2f, triggered at %s",
		a.Value, a.When.Format(time.RFC3339),
	)
}
