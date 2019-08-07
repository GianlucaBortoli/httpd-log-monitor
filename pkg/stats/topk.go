package stats

import "github.com/wangjia184/sortedset"

type TopK struct {
	k         int
	sortedSet *sortedset.SortedSet
}

func New(k int) *TopK {
	return &TopK{
		k:         k,
		sortedSet: sortedset.New(),
	}
}

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

func (t *TopK) TopK() []string {
	var out []string
	for i := 0; i < t.k; i++ {
		max := t.sortedSet.PopMax()
		// Append key only on valid elements. PopMax returns nil if the SortedSet is empty
		if max != nil {
			out = append(out, max.Key())
		}
	}
	return out
}

func (t *TopK) addOrUpdate(key string, score int64) bool {
	return t.sortedSet.AddOrUpdate(key, (sortedset.SCORE)(score), score)
}
