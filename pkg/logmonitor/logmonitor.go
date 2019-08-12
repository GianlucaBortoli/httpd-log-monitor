package logmonitor

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/internal/logparser"
	"github.com/cog-qlik/httpd-log-monitor/internal/tailer"
	"github.com/cog-qlik/httpd-log-monitor/pkg/metrics/manager"
	"github.com/hpcloud/tail"
)

// Monitor scrapes log files and derives statistics from it
type Monitor struct {
	parser       *logparser.HTTPd
	tailer       *tailer.Tailer
	statsManager *manager.Manager
	log          *log.Logger
	quitChan     chan struct{}
	startTime    time.Time
}

// New creates a monitor
func New(fileName string, alertPeriod, statsPeriod time.Duration, k int, threshold float64) (*Monitor, error) {
	l := log.New(os.Stderr, "", log.LstdFlags)

	m, err := manager.New(alertPeriod, statsPeriod, k, threshold, l)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		parser:       logparser.New(),
		tailer:       tailer.New(fileName),
		statsManager: m,
		log:          l,
		quitChan:     make(chan struct{}),
		startTime:    time.Now(),
	}, nil
}

// Start starts tailing and processing the log file in a separate goroutine
func (m *Monitor) Start() error {
	lines, err := m.tailer.Start()
	if err != nil {
		return fmt.Errorf("monitor start error: %v", err)
	}

	go m.startParsingTail(lines)
	m.statsManager.Start()
	return nil
}

// Stop stops the tailer and the processing of new log lines
func (m *Monitor) Stop() error {
	if err := m.tailer.Stop(); err != nil {
		return err
	}
	close(m.quitChan)
	m.statsManager.Stop()
	return nil
}

// Wait blocks until the tailer goroutine is in a dead state. Returns the reason for its death.
// If the main process is abruptly killed, this function never returns and the tailer may leak inotify
// watches in the Linux kernel. See https://godoc.org/github.com/hpcloud/tail#Tail.Cleanup) for more
// information.
func (m *Monitor) Wait() error {
	return m.tailer.Wait()
}

// startParsingTail is the loop where every log line is parsed and processed
func (m *Monitor) startParsingTail(lines <-chan *tail.Line) {
	for {
		select {
		case l := <-lines:
			logLine, err := m.checkLine(l)
			if err != nil {
				m.log.Println("[ERROR]", err)
				continue
			}
			m.statsManager.ObserveSection(logLine.Section)
			m.statsManager.ObserveRequest()
		case <-m.quitChan:
			m.log.Println("[INFO] exiting monitor")
			return
		}
	}
}

// checkLine ensures the input line (coming directly from the tailer) respects the layout defined
// in https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format.
// It returns an error also in case log line contains a date preceding the time start time of the
// monitor. This allows the caller to skip both malformed and old log lines.
func (m *Monitor) checkLine(line *tail.Line) (*logparser.Line, error) {
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
