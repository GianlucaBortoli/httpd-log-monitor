package topk

import "fmt"

// Item represents a data point for the TopK statistic
type Item struct {
	Key   string
	Score int64
}

func (i *Item) String() string {
	return fmt.Sprintf("key:%s, score:%d", i.Key, i.Score)
}
