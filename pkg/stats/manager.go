package stats

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/stats/topk"
)

type Manager struct {
	ticker       *time.Ticker
	sectionsTopK *topk.TopK
	sectionsChan chan *topk.Item
	quitChan     chan struct{}
	startOnce    sync.Once
	stopOnce     sync.Once
	log          *log.Logger
	started      bool
}

func NewManager(period time.Duration, k int) *Manager {
	return &Manager{
		ticker:       time.NewTicker(period),
		sectionsTopK: topk.New(k),
		sectionsChan: make(chan *topk.Item),
		quitChan:     make(chan struct{}),
		log:          log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (m *Manager) Start() {
	m.startOnce.Do(func() {
		go m.loop()
		m.started = true
	})
}

func (m *Manager) Stop() {
	if !m.started {
		return
	}
	// Ensure signal on quitChan is sent only once
	m.stopOnce.Do(func() {
		m.quitChan <- struct{}{}
	})
}

func (m *Manager) ObserveSection(s string) {
	if !m.started {
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
	for _, s := range sections {
		m.log.Println(s.String())
	}
}
