package alert

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestAlert() *Alert {
	a, _ := New(50*time.Millisecond, 100*time.Millisecond, 1, nil)
	return a
}

func TestNew(t *testing.T) {
	a, err := New(time.Second, 2*time.Second, 100, nil)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.IsType(t, &Alert{}, a)
	assert.Equal(t, float64(100), a.threshold)
	assert.False(t, a.firing)
}

func TestNew_WrongAlertPeriod(t *testing.T) {
	a, err := New(time.Second, 0, 100, nil)
	assert.Error(t, err)
	assert.Nil(t, a)
}

func TestNew_WrongStatPeriod(t *testing.T) {
	a, err := New(0, time.Second, 100, nil)
	assert.Error(t, err)
	assert.Nil(t, a)
}

func TestAlert_IncrByWithAlert(t *testing.T) {
	a, _ := New(50*time.Millisecond, 100*time.Millisecond, 1, nil)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// An alert message is expected since increment exceeds the
		// alert threshold
		msg := <-a.Alerts
		fmt.Println(msg)
		assert.NotEmpty(t, msg)
	}()

	// This should send a an alert message in the channel since the new average
	// is over the threshold. The goroutine above ensures we get the message and
	// we wait for it to return so we are sure the message has been received correctly
	err := a.incrBy(10)
	assert.NoError(t, err)
	wg.Wait()
}

func TestAlert_IncrByNoAlert(t *testing.T) {
	a, _ := New(time.Second, time.Hour, 100, nil)
	// This should not send an alert message, so the function shouldn't block on
	// sending the message on te channel.
	// Hence, if the tests exits it means that no message has been sent in the channel
	// since there's nothing reading from it.
	err := a.incrBy(1)
	assert.NoError(t, err)
}

func TestAlert_IncrByErr(t *testing.T) {
	a := getTestAlert()
	err := a.incrBy(-100)
	assert.Error(t, err)
}

func TestAlert_Reset(t *testing.T) {
	a := getTestAlert()

	err := a.metric.IncrBy(10)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), a.metric.GetCount())

	a.reset()
	assert.Equal(t, float64(0), a.metric.GetCount())
}
