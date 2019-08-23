package rate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a, err := New(time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, float64(0), a.count)
	assert.Equal(t, time.Second, a.windowSize)
}

func TestNew_Err(t *testing.T) {
	a, err := New(0)
	assert.Error(t, err)
	assert.Nil(t, a)
}

func TestAvg_IncrBy(t *testing.T) {
	a, err := New(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), a.count)

	err = a.IncrBy(1)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), a.count)

	err = a.IncrBy(0)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), a.count)

	err = a.IncrBy(-1)
	assert.Error(t, err)
	assert.Equal(t, float64(1), a.count)

	err = a.IncrBy(100)
	assert.NoError(t, err)
	assert.Equal(t, float64(101), a.count)
}

func TestAvg_Reset(t *testing.T) {
	a, err := New(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), a.count)

	err = a.IncrBy(10)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), a.count)

	a.Reset()
	assert.Equal(t, float64(0), a.count)
}

func TestAvg_AvgPerSec(t *testing.T) {
	a, err := New(2 * time.Minute)
	assert.NoError(t, err)

	// Add 100 data samples
	err = a.IncrBy(100)
	assert.NoError(t, err)

	// 2min = 120sec
	// 100samples / 120seconds = 0.8333333333333334req/sec
	avgPerSec := a.AvgPerSec()
	assert.Equal(t, float64(0.8333333333333334), avgPerSec)
}
