package stats

import (
	"fmt"

	"github.com/wangjia184/sortedset"
)

// TopK is an efficient data structure to store a scoreboard
type TopK struct {
	k         int
	sortedSet *sortedset.SortedSet
}

type TopKItem struct {
	key   string
	score int64
}

func (i *TopKItem) String() string {
	return fmt.Sprintf("key:%s, score:%d", i.key, i.score)
}

// NewTopK returns a new TopK metric
func NewTopK(k int) *TopK {
	return &TopK{
		k:         k,
		sortedSet: sortedset.New(),
	}
}

// IncrBy increments the score of item with key "key" of "incr".
// If key doesn't exist in the SortedSet it creates a new item with key "key" and score "incr".
// Inspired to https://redis.io/commands/zincrby
func (t *TopK) IncrBy(key string, incr int64) bool {
	if incr <= 0 {
		return false
	}

	item := t.sortedSet.Remove(key)
	if item == nil {
		// Add new item
		return t.addOrUpdate(key, incr)
	}
	// Update existing item
	newScore := (int64)(item.Score()) + incr
	return t.addOrUpdate(key, newScore)
}

// TopK returns at maximum "t.k" keys from the SortedSet
func (t *TopK) TopK() []*TopKItem {
	var out []*TopKItem
	for i := 0; i < t.k; i++ {
		max := t.sortedSet.PopMax()
		// Append key only on valid elements. PopMax returns nil if the SortedSet is empty
		if max != nil {
			out = append(out, &TopKItem{
				key:   max.Key(),
				score: int64(max.Score()),
			})
		}
	}
	return out
}

// Reset replaces the SortedSet with an empty one.
// This means that this deletes any data that was inside.
func (t *TopK) Reset() {
	t.sortedSet = sortedset.New()
}

func (t *TopK) addOrUpdate(key string, score int64) bool {
	return t.sortedSet.AddOrUpdate(key, (sortedset.SCORE)(score), score)
}
