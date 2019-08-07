package stats

import (
	"log"
	"os"
	"time"
)

type Manager struct {
	ticker       *time.Ticker
	sectionsTopK *TopK
	sectionsChan chan *sectionsObs
	quitChan     chan struct{}
	log          *log.Logger
	started      bool
}

type sectionsObs struct {
	key  string
	incr int64
}

func NewManager(period time.Duration, k int) *Manager {
	return &Manager{
		ticker:       time.NewTicker(period),
		sectionsTopK: NewTopK(k),
		sectionsChan: make(chan *sectionsObs),
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
	m.sectionsChan <- &sectionsObs{key: s, incr: 1}
}

func (m *Manager) loop() {
	for {
		select {
		case <-m.ticker.C:
			m.printSections(m.sectionsTopK.TopK())
			m.sectionsTopK.Reset()
		case o := <-m.sectionsChan:
			if ok := m.sectionsTopK.IncrBy(o.key, o.incr); !ok {
				m.log.Printf("[ERROR] cannot incremet key %s by %d\n", o.key, o.incr)
			}
		case <-m.quitChan:
			m.log.Println("[INFO] exiting stats manager")
			close(m.sectionsChan)
			close(m.quitChan)
		}
	}
}

func (m *Manager) printSections(sections []*TopKItem) {
	for _, s := range sections {
		m.log.Println(s.String())
	}
}
