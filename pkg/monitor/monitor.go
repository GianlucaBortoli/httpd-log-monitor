package monitor

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/internal/logparser"
	"github.com/cog-qlik/httpd-log-monitor/internal/tailer"
	"github.com/hpcloud/tail"
)

// Monitor scrapes log files and derives statistics from it
type Monitor struct {
	parser    *logparser.HTTPd
	tailer    *tailer.Tailer
	log       *log.Logger
	quitChan  chan struct{}
	startTime time.Time
}

// New creates a monitor. Returns the monitor and an error
func New(tailer *tailer.Tailer) (*Monitor, error) {
	if tailer == nil {
		return nil, fmt.Errorf("monitor needs a log tailer")
	}
	return &Monitor{
		parser:    logparser.New(),
		tailer:    tailer,
		log:       log.New(os.Stderr, "", log.LstdFlags),
		quitChan:  make(chan struct{}),
		startTime: time.Now(),
	}, nil
}

// Start starts tailing and processing the log file in a separate goroutine
func (m *Monitor) Start() error {
	lines, err := m.tailer.Start()
	if err != nil {
		return fmt.Errorf("monitor start error: %v", err)
	}

	go m.startParsingTail(lines)
	return nil
}

// Stop stops the tailer and the processing of new log lines
func (m *Monitor) Stop() error {
	if err := m.tailer.Stop(); err != nil {
		return err
	}
	m.quitChan <- struct{}{}
	return nil
}

// Wait blocks until the tailer goroutine is in a dead state.
// Returns the reason for its death.
func (m *Monitor) Wait() error {
	return m.tailer.Wait()
}

// startParsingTail is the loop where every log line is parsed and processed
func (m *Monitor) startParsingTail(lines <-chan *tail.Line) {
	for {
		select {
		case l := <-lines:
			good, err := m.filterLine(l)
			if err != nil {
				m.log.Println("[ERROR]", err)
			}
			m.log.Println(good)
		case <-m.quitChan:
			m.log.Println("[INFO] exiting monitor")
			return
		}
	}
}

func (m *Monitor) filterLine(line *tail.Line) (*logparser.Line, error) {
	if line == nil {
		return nil, fmt.Errorf("nil line")
	}
	parsedLine, err := m.parser.ParseLine(line.Text)
	if err != nil {
		return nil, fmt.Errorf("error parsing line: %v", err)
	}
	// Skip log lines whose date is before the start of the monitor.
	// This avoids to consider stale data for any later usage (eg. stats)
	if m.isOldLine(parsedLine) {
		return nil, fmt.Errorf("old log line detected. log time: %s, monitor start time: %s",
			parsedLine.Date.String(), m.startTime.String())
	}
	return parsedLine, nil
}

func (m *Monitor) isOldLine(line *logparser.Line) bool {
	if line == nil {
		return true
	}
	if line.Date.Before(m.startTime) {
		return true
	}
	return false
}
