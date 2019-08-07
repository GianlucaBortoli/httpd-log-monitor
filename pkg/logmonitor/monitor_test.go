package logmonitor

import (
	"os"
	"testing"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/internal/logparser"
	"github.com/cog-qlik/httpd-log-monitor/internal/testutils"
	"github.com/hpcloud/tail"
	"github.com/stretchr/testify/assert"
)

func getTestMonitor() (*Monitor, *os.File) {
	f, _ := testutils.CreateTestFile()
	m := New(f.Name())
	return m, f
}

func TestNew(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	m := New(f.Name())
	assert.NotNil(t, m)
	assert.IsType(t, &Monitor{}, m)
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

func TestMonitor_FilterLine(t *testing.T) {
	layout := "02/Jan/2006:15:04:05 -0700"
	pastDate := "09/May/2018:16:00:39 +0000"
	futureDate := "09/May/2099:16:00:39 +0000"
	futureDateTime, _ := time.Parse(layout, futureDate)

	testCases := []struct {
		line          *tail.Line
		expParsedLine *logparser.Line
		expErr        bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&tail.Line{},
			nil,
			true,
		},
		{
			// Date in the past
			&tail.Line{
				Text: `127.0.0.1 asd james [` + pastDate + `] "GET /report HTTP/1.0" 200 123`,
				Err:  nil,
				Time: time.Now(),
			},
			nil,
			true,
		},
		{
			// Date in the future
			&tail.Line{
				Text: `127.0.0.1 asd james [` + futureDate + `] "GET /report HTTP/1.0" 200 123`,
				Err:  nil,
				Time: time.Now(),
			},
			&logparser.Line{
				RemoteHost:    "127.0.0.1",
				RemoteLogName: "asd",
				User:          "james",
				Date:          futureDateTime,
				Method:        "GET",
				Section:       "/report",
				Protocol:      "HTTP/1.0",
				StatusCode:    200,
				ContentLength: 123,
			},
			false,
		},
	}

	m, _ := getTestMonitor()

	for _, tt := range testCases {
		parsed, err := m.filterLine(tt.line)
		assert.Equal(t, tt.expErr, err != nil)
		assert.Equal(t, tt.expParsedLine, parsed)
	}
}

func TestMonitor_IsOldLine(t *testing.T) {
	testCases := []struct {
		line     *logparser.Line
		expIsOld bool
	}{
		{
			&logparser.Line{Date: time.Now().AddDate(0, 0, -1)},
			true,
		},
		{
			&logparser.Line{Date: time.Now().AddDate(0, 0, +1)},
			false,
		},
		{
			nil,
			true,
		},
	}

	m, _ := getTestMonitor()

	for _, tt := range testCases {
		isOld := m.isOldLine(tt.line)
		assert.Equal(t, tt.expIsOld, isOld)
	}
}
