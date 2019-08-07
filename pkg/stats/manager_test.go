package stats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestManager() *Manager {
	return NewManager(50*time.Millisecond, 10)
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

	cnt := m.sectionsTopK.sortedSet.GetCount()
	assert.Equal(t, 0, cnt)

	m.ObserveSection("/foo")
	m.ObserveSection("/bar")
	m.ObserveSection("/baz")

	time.Sleep(100 * time.Millisecond) // give ticker the time to fire
}
