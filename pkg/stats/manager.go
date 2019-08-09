package stats

import (
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/stats/topk"
)

// Manager manages all the statistics computed from logs
type Manager struct {
	ticker       *time.Ticker
	sectionsTopK *topk.TopK
	sectionsChan chan *topk.Item
	quitChan     chan struct{}
	startOnce    sync.Once
	stopOnce     sync.Once
	log          *log.Logger
	started      int32 // 0 stopped, 1 started
}

// NewManager returns a new manager
func NewManager(period time.Duration, k int, l *log.Logger) *Manager {
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}

	return &Manager{
		ticker:       time.NewTicker(period),
		sectionsTopK: topk.New(k),
		sectionsChan: make(chan *topk.Item),
		quitChan:     make(chan struct{}),
		log:          l,
	}
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

func (m *Manager) loop() {
	for {
		select {
		case <-m.ticker.C:
			m.printSections(m.sectionsTopK.TopK())
			m.sectionsTopK.Reset()
		case i := <-m.sectionsChan:
			if ok := m.sectionsTopK.IncrBy(i); !ok {
				m.log.Printf("[ERROR] cannot incremet key %s by %d\n", i.Key, i.Score)
			}
		case <-m.quitChan:
			m.log.Println("[INFO] exiting stats manager")
			return
		}
	}
}

func (m *Manager) printSections(sections []*topk.Item) {
	if len(sections) == 0 {
		m.log.Println("no sections in the last period")
		return
	}
	m.log.Println("------------------------------")
	for _, s := range sections {
		m.log.Println(s.String())
	}
}
