package aggregator

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Aggregator interface {
	Max(from, to time.Time) (*AggregationResult, error)
	Median(from, to time.Time) (*AggregationResult, error)
}

type AggregationResult struct {
	Distance int64
	Duration time.Duration

	DistanceByBucket int64
	DurationByBucket time.Duration
}

// for home-task purposes
func NewInMemory[T Row](inputCh <-chan T, bucketSize BucketSize) (Aggregator, error) {
	if inputCh == nil {
		return nil, fmt.Errorf("sourceCh cannot be nil")
	}

	a := aggregator[T]{
		mu:      &sync.RWMutex{},
		buckets: make(map[string]*list.List),
	}
	a.bucketKeyFn, a.nextKeyFn = bucketFuncsFromSize(bucketSize)

	go func() {
		for row := range inputCh {
			a.insertRow(row)
		}
	}()
	return &a, nil
}

func (a *aggregator[T]) insertRow(row T) {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := a.bucketKeyFn(row.GetTimestamp())

	var bucketElems *list.List

	var ok bool
	if bucketElems, ok = a.buckets[key]; !ok {
		bucketElems = &list.List{}
		bucketElems.Init()
		a.buckets[key] = bucketElems
	}

	e := bucketElems.PushFront(row)
	for e.Next() != nil {
		// we store each bucket as a list of RunSessions ordered by time
		if e.Next().Value.(T).GetTimestamp().After(row.GetTimestamp()) {
			bucketElems.MoveAfter(e, e.Next())
			continue
		}
		break
	}
}

type aggregator[T Row] struct {
	mu      *sync.RWMutex
	buckets map[string]*list.List

	bucketKeyFn func(time.Time) string
	nextKeyFn   func(string) string
}
