package aggregator

import (
	"github.com/zyedidia/generic/list"

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
		buckets: make(map[string]*list.List[T]),
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

	var bucketElems *list.List[T]

	var ok bool
	if bucketElems, ok = a.buckets[key]; !ok {
		bucketElems = list.New[T]()
		a.buckets[key] = bucketElems
	}

	var inserted bool
	for ptr := bucketElems.Front; ptr != nil && !inserted; ptr = ptr.Next {
		if ptr.Value.GetTimestamp().After(row.GetTimestamp()) {
			// if current node has older ts than inserted then keep traversing
			continue
		}
		bucketElems.InsertBefore(ptr, &list.Node[T]{Value: row})
		inserted = true
	}
	if !inserted {
		bucketElems.PushFront(row)
	}
}

type aggregator[T Row] struct {
	mu      *sync.RWMutex
	buckets map[string]*list.List[T]

	bucketKeyFn func(time.Time) string
	nextKeyFn   func(string) string
}
