package stats

import (
	"testing"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/stats/topk"
	"github.com/stretchr/testify/assert"
)

func getTestManager() *Manager {
	return NewManager(50*time.Millisecond, 10, nil)
}

func TestNewManager(t *testing.T) {
	m := getTestManager()
	assert.NotNil(t, m)
	assert.IsType(t, &Manager{}, m)
}

func TestManager_ObserveSectionNotStarted(t *testing.T) {
	m := getTestManager()
	m.ObserveSection("/foo")
}

func TestManager_ObserveSection(t *testing.T) {
	m := getTestManager()
	m.Start()

	cnt := m.sectionsTopK.GetCount()
	assert.Equal(t, 0, cnt)

	m.ObserveSection("/foo")
	m.ObserveSection("/foo")
	m.ObserveSection("/bar")
	m.ObserveSection("/bar")
	m.ObserveSection("/bar")
	m.ObserveSection("/baz")
	// Give ticker some time to fire so I see console output
	time.Sleep(70 * time.Millisecond)
}

func TestManager_Start(t *testing.T) {
	m := getTestManager()
	m.Start()
}

func TestManager_StartMultiple(t *testing.T) {
	m := getTestManager()
	m.Start()
	m.Start()
	m.Start()
	m.Start()
}

func TestManager_StartAndStop(t *testing.T) {
	m := getTestManager()
	m.Start()
	m.Stop()
}

func TestManager_StartAndMultipleStop(t *testing.T) {
	m := getTestManager()
	m.Start()
	m.Stop()
	m.Stop()
	m.Stop()
}

func TestManager_StopNotStarted(t *testing.T) {
	m := getTestManager()
	m.Stop()
}

func TestManager_printSections(t *testing.T) {
	m := getTestManager()
	m.printSections(nil)
	m.printSections([]*topk.Item{{"a", 1}, {"b", 2}})
}
