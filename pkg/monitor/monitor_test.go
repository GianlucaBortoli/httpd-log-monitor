package monitor

import (
	"os"
	"testing"

	"github.com/cog-qlik/httpd-log-monitor/internal/tailer"
	"github.com/cog-qlik/httpd-log-monitor/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func getTestMonitor() (*Monitor, *os.File) {
	f, _ := testutils.CreateTestFile()
	ta := tailer.New(f.Name())
	m, _ := New(ta)
	return m, f
}

func TestNew(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	ta := tailer.New(f.Name())
	m, err := New(ta)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestNew_NilTailerErr(t *testing.T) {
	m, err := New(nil)
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestMonitor_Start(t *testing.T) {
	m, f := getTestMonitor()
	defer testutils.RemoveTestFile(f)

	err := m.Start()
	assert.NoError(t, err)
}

func TestMonitor_StartWhenAlreadyStarted(t *testing.T) {
	m, f := getTestMonitor()
	defer testutils.RemoveTestFile(f)

	err := m.Start()
	assert.NoError(t, err)

	err2 := m.Start()
	assert.Error(t, err2)
}

func TestMonitor_StartAndStop(t *testing.T) {
	m, f := getTestMonitor()
	defer testutils.RemoveTestFile(f)

	err := m.Start()
	assert.NoError(t, err)

	err2 := m.Stop()
	assert.NoError(t, err2)
}

func TestMonitor_StopWhenNotStarted(t *testing.T) {
	m, f := getTestMonitor()
	defer testutils.RemoveTestFile(f)

	err := m.Stop()
	assert.Error(t, err)
}

func TestMonitor_StartStopWait(t *testing.T) {
	m, f := getTestMonitor()
	defer testutils.RemoveTestFile(f)

	err := m.Start()
	assert.NoError(t, err)

	err2 := m.Stop()
	assert.NoError(t, err2)

	err3 := m.Wait()
	assert.NoError(t, err3)
}

func TestMonitor_WaitWhenNotStarted(t *testing.T) {
	m, f := getTestMonitor()
	defer testutils.RemoveTestFile(f)

	err := m.Wait()
	assert.Error(t, err)
}
