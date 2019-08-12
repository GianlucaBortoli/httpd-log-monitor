package manager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestManager() *Manager {
	m, _ := New(50*time.Millisecond, 50*time.Millisecond, 10, 10, nil)
	return m
}

func TestNewManager(t *testing.T) {
	m := getTestManager()
	assert.NotNil(t, m)
	assert.IsType(t, &Manager{}, m)
}

func TestManager_ObserveSection(t *testing.T) {
	m := getTestManager()
	m.Start()

	cnt := m.sectionsTopK.Count()
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

func TestManager_ObserveSectionNotStarted(t *testing.T) {
	m := getTestManager()
	m.ObserveSection("/foo")
}

func TestManager_ObserveUser(t *testing.T) {
	m := getTestManager()
	m.Start()

	cnt := m.sectionsTopK.Count()
	assert.Equal(t, 0, cnt)

	m.ObserveUser("1")
	m.ObserveUser("2")
	m.ObserveUser("3")
	m.ObserveUser("2")
	m.ObserveUser("2")
	// Give ticker some time to fire so I see console output
	time.Sleep(70 * time.Millisecond)
}

func TestManager_ObserveUserNotStarted(t *testing.T) {
	m := getTestManager()
	m.ObserveUser("1")
}

func TestManager_ObserveRequest(t *testing.T) {
	m := getTestManager()
	m.Start()

	cnt := m.sectionsTopK.Count()
	assert.Equal(t, 0, cnt)

	m.ObserveRequest()
	m.ObserveRequest()
	m.ObserveRequest()
	m.ObserveRequest()
	// Give ticker some time to fire so I see console output
	time.Sleep(70 * time.Millisecond)
}

func TestManager_ObserveRequestNotStarted(t *testing.T) {
	m := getTestManager()
	m.ObserveRequest()
}

func TestManager_ObserveStatusCode(t *testing.T) {
	m := getTestManager()
	m.Start()

	cnt := m.sectionsTopK.Count()
	assert.Equal(t, 0, cnt)

	m.ObserveStatusCode(200)
	m.ObserveStatusCode(300)
	m.ObserveStatusCode(400)
	m.ObserveStatusCode(500)
	// Give ticker some time to fire so I see console output
	time.Sleep(70 * time.Millisecond)
}

func TestManager_ObserveStatusCodeNotStarted(t *testing.T) {
	m := getTestManager()
	m.ObserveStatusCode(1)
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

func TestManager_printTopK(t *testing.T) {
	m := getTestManager()
	m.printTopK(m.sectionsTopK)
}

func TestIsErrorStatusCode(t *testing.T) {
	testCases := []struct {
		statusCode int
		expIsError bool
	}{
		{100, true},
		{199, true},
		{200, false},
		{300, false},
		{399, false},
		{400, true},
		{500, true},
	}

	for _, tt := range testCases {
		isErr := isErrorStatusCode(tt.statusCode)
		assert.Equal(t, tt.expIsError, isErr)
	}
}
