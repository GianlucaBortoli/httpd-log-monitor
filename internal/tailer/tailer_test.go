package tailer

import (
	"fmt"
	"sync"
	"testing"

	"github.com/cog-qlik/httpd-log-monitor/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tailer := New("asd")
	assert.NotNil(t, tailer)
	assert.IsType(t, &Tailer{}, tailer)
}

func TestTailer_Start(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	tailer := New(f.Name())
	var wg sync.WaitGroup
	linesToWrite := 10
	linesCnt := 0

	// Appender
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer tailer.Stop() // Stop tailer when the last line is written

		for i := 0; i < linesToWrite; i++ {
			_, err := f.WriteString(fmt.Sprintf("Line %d\n", i))
			assert.NoError(t, err)
			syncErr := f.Sync() // Make sure every line is written to disk
			assert.NoError(t, syncErr)
		}
	}()

	// Tailer
	wg.Add(1)
	go func() {
		defer wg.Done()

		lines, err := tailer.Start()
		assert.NoError(t, err)
		for l := range lines {
			linesCnt += 1
			fmt.Println(l.Text)
			assert.NotEmpty(t, l)
		}
	}()

	wg.Wait()
	assert.Equal(t, linesToWrite, linesCnt)
}

func TestTailer_StartAlreadyStarted(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	tailer := New(f.Name())
	lines, err := tailer.Start()
	assert.NoError(t, err)
	assert.NotNil(t, lines)

	lines2, err2 := tailer.Start()
	assert.Error(t, err2)
	assert.Nil(t, lines2)
}

func TestTailer_StartFileNotExistsErr(t *testing.T) {
	tailer := New("asd")
	lines, err := tailer.Start()
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestTailer_StartAndStop(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	tailer := New(f.Name())
	lines, startErr := tailer.Start()
	assert.NoError(t, startErr)
	assert.NotNil(t, lines)

	stopErr := tailer.Stop()
	assert.NoError(t, stopErr)
}

func TestTailer_StopWhenNotStarted(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	tailer := New(f.Name())
	stopErr := tailer.Stop()
	assert.Error(t, stopErr)
}

func TestTailer_WaitAfterStartAndStop(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	tailer := New(f.Name())
	lines, startErr := tailer.Start()
	assert.NoError(t, startErr)
	assert.NotNil(t, lines)

	stopErr := tailer.Stop()
	assert.NoError(t, stopErr)

	waitErr := tailer.Wait()
	assert.NoError(t, waitErr)
}

func TestTailer_WaitNotStarted(t *testing.T) {
	f, err := testutils.CreateTestFile()
	assert.NoError(t, err)
	defer testutils.RemoveTestFile(f)

	tailer := New(f.Name())
	waitErr := tailer.Wait()
	assert.Error(t, waitErr)
}
