package metrics

import (
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/secavg"
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
	reqSec     *secavg.SecAvg
	reqSecChan chan float64
}

// NewManager returns a new manager
func NewManager(statsPeriod time.Duration, k int, l *log.Logger) (*Manager, error) {
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	r, err := secavg.New(statsPeriod)
	if err != nil {
		return nil, err
	}

	return &Manager{
		metricsTicker: time.NewTicker(statsPeriod),
		quitChan:      make(chan struct{}),
		log:           l,
		sectionsTopK:  topk.New(k),
		sectionsChan:  make(chan *topk.Item),
		reqSec:        r,
		reqSecChan:    make(chan float64),
	}, nil
}

// Start starts the stats manager
func (m *Manager) Start() {
	m.startOnce.Do(func() {
		go m.loop()
		atomic.StoreInt32(&m.started, 1)
	})
}

// Stop stops the stats manager
func (m *Manager) Stop() {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}
	// Ensure signal on quitChan is sent only once
	m.stopOnce.Do(func() {
		m.quitChan <- struct{}{}
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
		case <-m.quitChan:
			m.log.Println("[INFO] exiting stats manager")
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
