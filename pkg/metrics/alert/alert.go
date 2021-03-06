package alert

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GianlucaBortoli/httpd-log-monitor/pkg/metrics/rate"
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
	Alerts    chan *msg // Alerts are sent here
}

// New returns the alert manager with the specified alerting period and threshold
func New(period time.Duration, threshold float64, l *log.Logger) (*Alert, error) {
	if period == 0 {
		return nil, fmt.Errorf("cannot create alert with time window of width 0")
	}
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	m, err := rate.New(period)
	if err != nil {
		return nil, fmt.Errorf("cannot create rate metric for alert: %v", err)
	}

	return &Alert{
		ticker:    time.NewTicker(period),
		log:       l,
		metric:    m,
		threshold: threshold,
		incrChan:  make(chan float64),
		quitChan:  make(chan struct{}),
		Alerts:    make(chan *msg, 100),
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
			a.metric.Reset()
		case msg := <-a.Alerts:
			a.print(msg)
		case i := <-a.incrChan:
			if err := a.metric.IncrBy(i); err != nil {
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
	avg := a.metric.AvgPerSec()

	if !a.firing && avg >= a.threshold {
		a.Alerts <- &msg{
			Type:  highTraffic,
			Value: avg,
			When:  time.Now(),
		}
		a.firing = true // alert firing
	}
	if a.firing && avg < a.threshold {
		a.Alerts <- &msg{
			Type:  resolved,
			Value: avg,
			When:  time.Now(),
		}
		a.firing = false // alert resolved
	}
}

func (a *Alert) print(msg *msg) {
	switch msg.Type {
	case highTraffic:
		a.log.Println("[ALERT]", msg.String())
	case resolved:
		a.log.Println("[RESOLVED]", msg.String())
	default:
		a.log.Println(msg.String())
	}
}
