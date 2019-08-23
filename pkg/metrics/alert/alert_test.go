package alert

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestAlert() *Alert {
	a, _ := New(100*time.Millisecond, 1, nil)
	return a
}

func TestNew(t *testing.T) {
	a, err := New(time.Second, 100, nil)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, float64(100), a.threshold)
	assert.False(t, a.firing)
}

func TestNew_WrongPeriod(t *testing.T) {
	a, err := New(0, 100, nil)
	assert.Error(t, err)
	assert.Nil(t, a)
}

func TestAlert_IncrByNoAlert(t *testing.T) {
	a, _ := New(time.Hour, 100, nil)
	// This should not send an alert message, so the function shouldn't block on
	// sending the message on te channel.
	// Hence, if the tests exits it means that no message has been sent in the channel
	// since there's nothing reading from it.
	err := a.metric.IncrBy(1)
	assert.NoError(t, err)
}

func TestAlert_Reset(t *testing.T) {
	a := getTestAlert()

	err := a.metric.IncrBy(10)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), a.metric.Count())

	a.metric.Reset()
	assert.Equal(t, float64(0), a.metric.Count())
}

func TestAlert_Start(t *testing.T) {
	a := getTestAlert()
	a.Start()
	assert.True(t, a.started)
}

func TestAlert_StartMultiple(t *testing.T) {
	a := getTestAlert()
	a.Start()
	a.Start()
	a.Start()
	a.Start()
	a.Start()
	assert.True(t, a.started)
}

func TestAlert_Stop(t *testing.T) {
	a := getTestAlert()
	a.Stop()
}

func TestAlert_StartAndStop(t *testing.T) {
	a := getTestAlert()
	a.Start()
	assert.True(t, a.started)
	a.Stop()
	assert.False(t, a.started)
}

func TestAlert_IncrByStarted(t *testing.T) {
	a := getTestAlert()
	a.Start()
	a.IncrBy(1)
}

func TestAlert_IncrByStopped(t *testing.T) {
	a := getTestAlert()
	a.IncrBy(1)
}

func TestAlert_checkThresholdWithAlerts(t *testing.T) {
	a, _ := New(5*time.Second, 1.9, nil)
	a.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 10samples / 5sec = 2req/sec
		// Threshold is set to 1.9 so I should get an alert message
		msg := <-a.Alerts

		fmt.Println(msg)
		assert.NotNil(t, msg)
		assert.Equal(t, highTraffic, msg.Type)

		// After another period with no traffic, I should get a resolved message
		msg = <-a.Alerts
		fmt.Println(msg)
		assert.NotNil(t, msg)
		assert.Equal(t, resolved, msg.Type)
	}()

	// Add 10 data samples only once to simulate an initial spike of requests
	a.IncrBy(10)
	wg.Wait()
}
