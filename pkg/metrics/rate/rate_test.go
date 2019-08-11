package secavg

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
	assert.IsType(t, &Rate{}, a)
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

func TestAvg_GetAvgPerSec(t *testing.T) {
	a, err := New(10 * time.Minute)
	assert.NoError(t, err)

	err = a.IncrBy(60)
	assert.NoError(t, err)
	avgPerSec := a.GetAvgPerSec()
	assert.Equal(t, 0.1, avgPerSec)
}
