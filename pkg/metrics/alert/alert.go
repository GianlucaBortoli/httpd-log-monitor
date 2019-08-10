package alert

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/secavg"
)

type Alert struct {
	ticker    *time.Ticker
	log       *log.Logger
	metric    *secavg.SecAvg
	threshold float64
	firing    bool
	started   bool
	incrChan  chan float64
	quitChan  chan struct{}
	Alerts    chan *Msg
}

func New(statPeriod, alertPeriod time.Duration, threshold float64, l *log.Logger) (*Alert, error) {
	if alertPeriod == 0 {
		return nil, fmt.Errorf("cannot create alert with period %d", alertPeriod)
	}
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	m, err := secavg.New(statPeriod)
	if err != nil {
		return nil, fmt.Errorf("cannot create alert: %v", err)
	}

	return &Alert{
		ticker:    time.NewTicker(alertPeriod),
		log:       l,
		metric:    m,
		threshold: threshold,
		incrChan:  make(chan float64),
		quitChan:  make(chan struct{}),
		Alerts:    make(chan *Msg, 100),
	}, nil
}

func (a *Alert) Start() {
	if a.started {
		return
	}
	go a.loop()
	a.started = true
}

func (a *Alert) Stop() {
	if !a.started {
		return
	}
	a.quitChan <- struct{}{}
	a.started = false
}

func (a *Alert) IncrBy(i float64) {
	if !a.started {
		return
	}
	a.incrChan <- i
}

func (a *Alert) loop() {
	for {
		select {
		case <-a.ticker.C:
			a.checkThreshold()
			a.reset()
		case msg := <-a.Alerts:
			a.printAlert(msg)
		case i := <-a.incrChan:
			if err := a.incrBy(i); err != nil {
				a.log.Println("[ERROR]", err)
			}
		case <-a.quitChan:
			a.log.Println("[INFO] alert event loop exit")
		}
	}
}

func (a *Alert) incrBy(i float64) error {
	return a.metric.IncrBy(i)
}

func (a *Alert) checkThreshold() {
	avg := a.metric.GetAvgPerSec()

	if avg >= a.threshold {
		a.Alerts <- &Msg{
			Type:  HighTraffic,
			Value: avg,
			When:  time.Now(),
		}
		a.firing = true
	}
	if a.firing && avg < a.threshold {
		a.Alerts <- &Msg{
			Type:  Resolved,
			Value: avg,
			When:  time.Now(),
		}
		a.firing = false
	}
}

func (a *Alert) reset() {
	a.metric.Reset()
}

func (a *Alert) printAlert(msg *Msg) {
	switch msg.Type {
	case HighTraffic:
		a.log.Println("[ALERT]", msg.String())
	case Resolved:
		a.log.Println("[RESOLVED]", msg.String())
	default:
		a.log.Println(msg.String())
	}
}
