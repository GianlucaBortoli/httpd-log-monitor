package monitor

import (
	"fmt"
	"log"
	"os"

	"github.com/cog-qlik/httpd-log-monitor/internal/logparser"
	"github.com/cog-qlik/httpd-log-monitor/internal/tailer"
	"github.com/hpcloud/tail"
)

type Monitor struct {
	parser   *logparser.HTTPd
	tailer   *tailer.Tailer
	log      *log.Logger
	quitChan chan struct{}
}

func New(tailer *tailer.Tailer) (*Monitor, error) {
	if tailer == nil {
		return nil, fmt.Errorf("monitor needs a log tailer")
	}
	return &Monitor{
		parser:   logparser.New(),
		tailer:   tailer,
		log:      log.New(os.Stderr, "", log.LstdFlags),
		quitChan: make(chan struct{}),
	}, nil
}

func (m *Monitor) Start() error {
	lines, err := m.tailer.Start()
	if err != nil {
		return fmt.Errorf("monitor start error: %v", err)
	}

	go m.startParsingTail(lines)
	return nil
}

func (m *Monitor) Stop() error {
	if err := m.tailer.Stop(); err != nil {
		return err
	}
	m.quitChan <- struct{}{}
	return nil
}

func (m *Monitor) Wait() error {
	return m.tailer.Wait()
}

func (m *Monitor) startParsingTail(lines <-chan *tail.Line) {
	for {
		select {
		case l := <-lines:
			m.handleLine(l)
		case <-m.quitChan:
			m.log.Println("[INFO] exiting monitor")
			return
		}
	}
}

func (m *Monitor) handleLine(line *tail.Line) {
	if line == nil {
		return
	}
	parsedLine, err := m.parser.ParseLine(line.Text)
	if err != nil {
		m.log.Printf("[ERROR] cannot parse line: %v\n", err)
	}
	m.log.Printf("%v", parsedLine)
}
