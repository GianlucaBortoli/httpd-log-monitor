package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangjia184/sortedset"
)

func TestNew(t *testing.T) {
	s := NewTopK(10)
	assert.NotNil(t, s)
	assert.IsType(t, &TopK{}, s)
}

func TestTopK_AddOrUpdate(t *testing.T) {
	s := NewTopK(10)
	assert.NotNil(t, s)

	ok := s.addOrUpdate("a", 1)
	assert.True(t, ok)
	ok = s.addOrUpdate("b", 2)
	assert.True(t, ok)
	ok = s.addOrUpdate("c", 3)
	assert.True(t, ok)
}

func TestTopK_TopKExact(t *testing.T) {
	s := NewTopK(3)
	assert.NotNil(t, s)

	ok := s.addOrUpdate("a", 1)
	assert.True(t, ok)
	ok = s.addOrUpdate("c", 3)
	assert.True(t, ok)
	ok = s.addOrUpdate("b", 2)
	assert.True(t, ok)

	val := s.TopK()
	assert.Len(t, val, 3)
	assert.Equal(t, val[0].key, "c")
	assert.Equal(t, val[1].key, "b")
	assert.Equal(t, val[2].key, "a")
}

func TestTopK_TopKLess(t *testing.T) {
	s := NewTopK(1)
	assert.NotNil(t, s)

	ok := s.addOrUpdate("a", 1)
	assert.True(t, ok)
	ok = s.addOrUpdate("c", 3)
	assert.True(t, ok)
	ok = s.addOrUpdate("b", 2)
	assert.True(t, ok)

	val := s.TopK()
	assert.Len(t, val, 1)
	assert.Equal(t, val[0].key, "c")
}

func TestTopK_TopKMore(t *testing.T) {
	s := NewTopK(10)
	assert.NotNil(t, s)

	ok := s.addOrUpdate("a", 1)
	assert.True(t, ok)
	ok = s.addOrUpdate("c", 3)
	assert.True(t, ok)
	ok = s.addOrUpdate("b", 2)
	assert.True(t, ok)

	val := s.TopK()
	assert.Len(t, val, 3)
	assert.Equal(t, val[0].key, "c")
	assert.Equal(t, val[1].key, "b")
	assert.Equal(t, val[2].key, "a")
}

func TestTopK_IncrBy(t *testing.T) {
	s := NewTopK(10)
	assert.NotNil(t, s)

	ok := s.IncrBy("a", 1)
	assert.True(t, ok)
	item := s.sortedSet.GetByKey("a")
	assert.Equal(t, (sortedset.SCORE)(1), item.Score())
	assert.Equal(t, "a", item.Key())

	ok = s.IncrBy("a", 1)
	assert.True(t, ok)
	item = s.sortedSet.GetByKey("a")
	assert.Equal(t, (sortedset.SCORE)(2), item.Score())
	assert.Equal(t, "a", item.Key())

	ok = s.IncrBy("a", 3)
	assert.True(t, ok)
	item = s.sortedSet.GetByKey("a")
	assert.Equal(t, (sortedset.SCORE)(5), item.Score())
	assert.Equal(t, "a", item.Key())
}

func TestTopK_IncrByErr(t *testing.T) {
	s := NewTopK(10)
	assert.NotNil(t, s)

	ok := s.IncrBy("a", -1)
	assert.False(t, ok)

	ok = s.IncrBy("a", 0)
	assert.False(t, ok)

	ok = s.IncrBy("a", 1)
	assert.True(t, ok)
}

func TestTopK_Reset(t *testing.T) {
	s := NewTopK(10)
	assert.NotNil(t, s)

	ok := s.IncrBy("a", 1)
	assert.True(t, ok)
	ok = s.IncrBy("b", 1)
	assert.True(t, ok)
	ok = s.IncrBy("c", 1)
	assert.True(t, ok)
	cnt := s.sortedSet.GetCount()
	assert.Equal(t, 3, cnt)

	s.Reset()
	cnt = s.sortedSet.GetCount()
	assert.Equal(t, 0, cnt)
}
