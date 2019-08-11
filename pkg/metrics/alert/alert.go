package alert

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/rate"
)

// Alert handles alerts for the per-second requests
type Alert struct {
	ticker    *time.Ticker
	log       *log.Logger
	metric    *rate.Rate
	threshold float64
	firing    bool
	started   bool
	incrChan  chan float64
	quitChan  chan struct{}
	Alerts    chan *Msg // Alerts are sent here
}

// New returns the alert manager with the specified metrics and alerting period and threshold
func New(statPeriod, alertPeriod time.Duration, threshold float64, l *log.Logger) (*Alert, error) {
	if alertPeriod == 0 {
		return nil, fmt.Errorf("cannot create alert with period %d", alertPeriod)
	}
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	m, err := rate.New(statPeriod)
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

// Start starts watching the requests per second metric
func (a *Alert) Start() {
	if a.started {
		return
	}
	go a.loop()
	a.started = true
}

// Stop stops the alert manager
func (a *Alert) Stop() {
	if !a.started {
		return
	}
	close(a.quitChan)
	a.started = false
}

// IncrBy observe a new incoming request
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

// checkThreshold checks whether the current requests per second average is above
// the threshold or not. Is also sends a message inside a.Alerts accordingly.
func (a *Alert) checkThreshold() {
	avg := a.metric.GetAvgPerSec()

	if !a.firing && avg >= a.threshold {
		a.Alerts <- &Msg{
			Type:  HighTraffic,
			Value: avg,
			When:  time.Now(),
		}
		a.firing = true // now the alert is firing
	}
	if a.firing && avg < a.threshold {
		a.Alerts <- &Msg{
			Type:  Resolved,
			Value: avg,
			When:  time.Now(),
		}
		a.firing = false // alert is resolved
	}
}

func (a *Alert) incrBy(i float64) error {
	return a.metric.IncrBy(i)
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
