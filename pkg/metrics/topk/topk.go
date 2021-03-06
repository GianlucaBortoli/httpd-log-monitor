package topk

import "github.com/wangjia184/sortedset"

// TopK is an efficient data structure to store a scoreboard
type TopK struct {
	k         int
	sortedSet *sortedset.SortedSet
}

// New returns a new TopK metric. Cannot return nil
func New(k int) *TopK {
	return &TopK{
		k:         k,
		sortedSet: sortedset.New(),
	}
}

// IncrBy increments the score of item with key "key" of "incr". If key doesn't exist in the
// SortedSet it creates a new item with key "key" and score "incr".
// The time complexity of this method is O(log({number of items in the SortedSet}))
// Inspired to https://redis.io/commands/zincrby.
func (t *TopK) IncrBy(i *Item) bool {
	if i == nil {
		return false
	}
	if i.Score <= 0 {
		return false
	}

	item := t.sortedSet.Remove(i.Key)
	if item == nil {
		// Add new item
		return t.addOrUpdate(i.Key, i.Score)
	}
	// Update existing item
	newScore := (int64)(item.Score()) + i.Score
	return t.addOrUpdate(i.Key, newScore)
}

// TopK returns at maximum "t.k" keys from the SortedSet
// The time complexity of this method is O(K * log({number of items in the SortedSet}))
func (t *TopK) TopK() []*Item {
	var out []*Item
	for i := 0; i < t.k; i++ {
		max := t.sortedSet.PopMax()
		// Append key only on valid elements. PopMax returns nil if the SortedSet is empty
		if max != nil {
			out = append(out, &Item{
				Key:   max.Key(),
				Score: int64(max.Score()),
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

// Count returns the number of items in the SortedSet.
// Used in unit tests only.
func (t *TopK) Count() int {
	return t.sortedSet.GetCount()
}

// addOrUpdate adds or updated key with score in the SortedSet.
// The time complexity of this method is O(log({number of items in the SortedSet}))
func (t *TopK) addOrUpdate(key string, score int64) bool {
	return t.sortedSet.AddOrUpdate(key, (sortedset.SCORE)(score), score)
}
