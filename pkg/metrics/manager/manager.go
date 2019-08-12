package manager

import (
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/alert"
	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/rate"
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
	// TopK status codes
	statusCodesTopK *topk.TopK
	statusCodesChan chan *topk.Item
	// TopK users
	usersTopK *topk.TopK
	usersChan chan *topk.Item
	// Req/sec metric
	reqSec     *rate.Rate
	reqSecChan chan float64
	// Err/sec metric
	errSec     *rate.Rate
	errSecChan chan float64
	// Req/sec alert
	reqSecAlert *alert.Alert
}

// New returns a new manager
func New(alertPeriod, statsPeriod time.Duration, k int, threshold float64, l *log.Logger) (*Manager, error) {
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	reqSec, mErr := rate.New(statsPeriod)
	if mErr != nil {
		return nil, mErr
	}

	errSec, eErr := rate.New(statsPeriod)
	if eErr != nil {
		return nil, eErr
	}

	a, aErr := alert.New(statsPeriod, alertPeriod, threshold, l)
	if aErr != nil {
		return nil, aErr
	}

	return &Manager{
		metricsTicker:   time.NewTicker(statsPeriod),
		quitChan:        make(chan struct{}),
		log:             l,
		sectionsTopK:    topk.New(k),
		sectionsChan:    make(chan *topk.Item),
		statusCodesTopK: topk.New(k),
		statusCodesChan: make(chan *topk.Item),
		usersTopK:       topk.New(k),
		usersChan:       make(chan *topk.Item),
		reqSec:          reqSec,
		reqSecChan:      make(chan float64),
		errSec:          errSec,
		errSecChan:      make(chan float64),
		reqSecAlert:     a,
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

// ObserveSection observes a data point for the sections TopK statistic
func (m *Manager) ObserveSection(s string) {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}
	m.sectionsChan <- &topk.Item{Key: s, Score: 1}
}

// ObserveUser observes a data point for the users TopK statistic
func (m *Manager) ObserveUser(s string) {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}
	m.usersChan <- &topk.Item{Key: s, Score: 1}
}

// ObserveRequest observes a data point for the requests per second statistic
func (m *Manager) ObserveRequest() {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}
	m.reqSecChan <- float64(1)
}

// ObserveStatusCode observes a data point for the errors per second statistic only if
// the input status code is considered an error
func (m *Manager) ObserveStatusCode(code int) {
	if atomic.LoadInt32(&m.started) == 0 {
		return
	}

	if isErrorStatusCode(code) {
		m.errSecChan <- float64(1)
	}
	m.statusCodesChan <- &topk.Item{Key: strconv.Itoa(code), Score: 1}
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
		case u := <-m.usersChan:
			if ok := m.usersTopK.IncrBy(u); !ok {
				m.log.Printf("[ERROR] cannot incremet key %s by %d\n", u.Key, u.Score)
			}
		case c := <-m.statusCodesChan:
			if ok := m.statusCodesTopK.IncrBy(c); !ok {
				m.log.Printf("[ERROR] cannot incremet key %s by %d\n", c.Key, c.Score)
			}
		case c := <-m.reqSecChan:
			if err := m.reqSec.IncrBy(c); err != nil {
				m.log.Println("[ERROR]", err)
			}
			m.reqSecAlert.IncrBy(c)
		case c := <-m.errSecChan:
			if err := m.errSec.IncrBy(c); err != nil {
				m.log.Println("[ERROR]", err)
			}
		case <-m.quitChan:
			m.log.Println("[INFO] exiting metrics manager event loop")
			return
		}
	}
}

func (m *Manager) printAllMetrics() {
	m.log.Println("------------------------------------------")
	m.printReqSec()
	m.printErrSec()
	m.log.Println("TopK sections:")
	m.printTopK(m.sectionsTopK)
	m.log.Println("TopK status codes:")
	m.printTopK(m.statusCodesTopK)
	m.log.Println("TopK users:")
	m.printTopK(m.usersTopK)
}

func (m *Manager) resetAllMetrics() {
	m.reqSec.Reset()
	m.errSec.Reset()
	m.sectionsTopK.Reset()
	m.statusCodesTopK.Reset()
	m.usersTopK.Reset()
}

func (m *Manager) printReqSec() {
	reqSec := m.reqSec.GetAvgPerSec()
	period := m.reqSec.GetWindowSize().String()
	m.log.Printf("%.2f req/s over last %s", reqSec, period)
}

func (m *Manager) printErrSec() {
	errSec := m.errSec.GetAvgPerSec()
	period := m.errSec.GetWindowSize().String()
	m.log.Printf("%.2f err/s over last %s", errSec, period)
}

func (m *Manager) printTopK(k *topk.TopK) {
	topK := k.TopK()
	if len(topK) == 0 {
		return
	}
	for _, s := range topK {
		m.log.Println(s.String())
	}
}

// isErrorStatusCode returns false if the status code is between 200 (included) and 400 (excluded)
// true otherwise
func isErrorStatusCode(code int) bool {
	if code >= 200 && code < 400 {
		return false
	}
	return true
}
