package stats

import (
	"log"
	"os"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/stats/topk"
)

type Manager struct {
	ticker       *time.Ticker
	sectionsTopK *topk.TopK
	sectionsChan chan *topk.Item
	quitChan     chan struct{}
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
	go m.loop()
	m.started = true
}

func (m *Manager) Stop() {
	m.quitChan <- struct{}{}
}

func (m *Manager) ObserveSection(s string) {
	if !m.started {
		return // noop
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
			close(m.sectionsChan)
			close(m.quitChan)
		}
	}
}

func (m *Manager) printSections(sections []*topk.Item) {
	for _, s := range sections {
		m.log.Println(s.String())
	}
}
