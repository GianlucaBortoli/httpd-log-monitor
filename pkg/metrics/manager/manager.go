package manager

import (
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/alert"
	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/topk"
)

// Manager manages all the statistics computed from logs
type Manager struct {
	metricsTicker *time.Ticker
	startOnce     sync.Once
	stopOnce      sync.Once
	log           *log.Logger
	started       int32 // 0 stopped, 1 started
	quitChan      chan struct{}
	// TopK sections metric
	sectionsTopK *topk.TopK
	sectionsChan chan *topk.Item
	// Req/sec metric
	reqSec     *rate.SecAvg
	reqSecChan chan float64
	// Req/sec alert
	reqSecAlert *alert.Alert
}

// New returns a new manager
func New(alertPeriod, statsPeriod time.Duration, k int, threshold float64, l *log.Logger) (*Manager, error) {
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	r, mErr := rate.New(statsPeriod)
	if mErr != nil {
		return nil, mErr
	}

	a, aErr := alert.New(statsPeriod, alertPeriod, threshold, l)
	if aErr != nil {
		return nil, aErr
	}

	return &Manager{
		metricsTicker: time.NewTicker(statsPeriod),
		quitChan:      make(chan struct{}),
		log:           l,
		sectionsTopK:  topk.New(k),
		sectionsChan:  make(chan *topk.Item),
		reqSec:        r,
		reqSecChan:    make(chan float64),
		reqSecAlert:   a,
	}, nil
}

// Start starts the stats manager
func (m *Manager) Start() {
	m.startOnce.Do(func() {
		go m.loop()
		m.reqSecAlert.Start()
		atomic.StoreInt32(&m.started, 1)
	})
}

// Stop stops the stats manager
func (m *Manager) Stop() {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}
	m.stopOnce.Do(func() {
		close(m.quitChan)
		m.reqSecAlert.Stop()
	})
}

// ObserveSection observe a data point for the sections TopK statistic
func (m *Manager) ObserveSection(s string) {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}
	m.sectionsChan <- &topk.Item{Key: s, Score: 1}
}

// ObserveRequest observe a data point for the requests per second statistic
func (m *Manager) ObserveRequest() {
	m.reqSecChan <- float64(1)
}

func (m *Manager) loop() {
	for {
		select {
		case <-m.metricsTicker.C:
			m.printAllMetrics()
			m.resetAllMetrics()
		case i := <-m.sectionsChan:
			if ok := m.sectionsTopK.IncrBy(i); !ok {
				m.log.Printf("[ERROR] cannot incremet key %s by %d\n", i.Key, i.Score)
			}
		case c := <-m.reqSecChan:
			if err := m.reqSec.IncrBy(c); err != nil {
				m.log.Println("[ERROR]", err)
			}
			m.reqSecAlert.IncrBy(c)
		case <-m.quitChan:
			m.log.Println("[INFO] exiting metrics manager event loop")
			return
		}
	}
}

func (m *Manager) printAllMetrics() {
	m.log.Println("------------------------------")
	m.printReqSec()
	m.printSections()
}

func (m *Manager) resetAllMetrics() {
	m.reqSec.Reset()
	m.sectionsTopK.Reset()
}

func (m *Manager) printReqSec() {
	reqPerSec := m.reqSec.GetAvgPerSec()
	period := m.reqSec.GetWindowSize().String()
	m.log.Printf("%.2f req/s over last %s", reqPerSec, period)
}

func (m *Manager) printSections() {
	sections := m.sectionsTopK.TopK()
	if len(sections) == 0 {
		m.log.Println("no sections in the last period")
		return
	}
	for _, s := range sections {
		m.log.Println(s.String())
	}
}
