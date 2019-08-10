package alert

import (
	"fmt"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/secavg"
)

type Alert struct {
	period    time.Duration
	metric    *secavg.SecAvg
	threshold float64
	Alerts    chan *alertMsg
}

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

func New(statPeriod, alertPeriod time.Duration, threshold float64) (*Alert, error) {
	if alertPeriod == 0 {
		return nil, fmt.Errorf("cannot create alert with period %d", alertPeriod)
	}

	m, err := secavg.New(statPeriod)
	if err != nil {
		return nil, fmt.Errorf("cannot create alert: %v", err)
	}

	return &Alert{
		period:    alertPeriod,
		metric:    m,
		threshold: threshold,
		Alerts:    make(chan *alertMsg),
	}, nil
}

func (a *Alert) IncrBy(i float64) error {
	err := a.metric.IncrBy(i)
	if err != nil {
		return fmt.Errorf("cannot increment metric for alert: %v", err)
	}

	if avg := a.getCurrAvgPerSec(i); avg >= a.threshold {
		a.Alerts <- &alertMsg{
			Value: avg,
			When:  time.Now(),
		}
	}
	return nil
}

func (a *Alert) Reset() {
	a.metric.Reset()
}

func (a *Alert) getCurrAvgPerSec(incr float64) float64 {
	oldCnt := a.metric.GetCount()
	oldAvg := a.metric.GetAvgPerSec()
	newCnt := oldCnt + incr
	newAvg := oldAvg / newCnt
	return newAvg
}
