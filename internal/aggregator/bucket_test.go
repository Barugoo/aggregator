package aggregator

import (
	"container/list"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Demonstrates the flexibility of bucket-based approach
*/
func TestForEachBucket(t *testing.T) {
	getKeyFn, nextKeyFn := bucketFuncsFromSize(HourBucketSize)

	agg := &aggregator[*RunSession]{
		mu:          &sync.RWMutex{},
		bucketKeyFn: getKeyFn,
		nextKeyFn:   nextKeyFn,
		buckets:     make(map[string]*list.List),
	}

	ts := time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC)
	// all the rows located in different buckets
	agg.insertRow(&RunSession{
		Timestamp: ts,
		Distance:  1,
	})
	agg.insertRow(&RunSession{
		Timestamp: ts.Add(time.Hour),
		Distance:  2,
	})
	agg.insertRow(&RunSession{
		Timestamp: ts.Add(time.Hour * 2),
		Distance:  3,
	})

	// normal case
	var bucketCount int64
	err := agg.forEachBucket(ts, ts.AddDate(1, 0, 0), func(bucketElems *list.List) error {
		bucketCount++
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(3), bucketCount)

	// error case
	err = agg.forEachBucket(ts, ts.AddDate(1, 0, 0), func(bucketElems *list.List) error {
		return errors.New("oops, something went wrong")
	})
	assert.Error(t, err)
}

func TestForEachBucketElem(t *testing.T) {
	getKeyFn, nextKeyFn := bucketFuncsFromSize(HourBucketSize)

	agg := &aggregator[*RunSession]{
		mu:          &sync.RWMutex{},
		bucketKeyFn: getKeyFn,
		nextKeyFn:   nextKeyFn,
		buckets:     make(map[string]*list.List),
	}

	ts := time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC)
	// all the rows located in one bucket
	agg.insertRow(&RunSession{
		Timestamp: ts,
		Distance:  1,
	})
	agg.insertRow(&RunSession{
		Timestamp: ts,
		Distance:  2,
	})
	agg.insertRow(&RunSession{
		Timestamp: ts,
		Distance:  3,
	})

	// normal case
	var distanceTotal int64
	err := agg.forEachBucketElem(agg.buckets[getKeyFn(ts)], func(elem *RunSession) error {
		distanceTotal += elem.GetDistance()
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(6), distanceTotal)

	// error case
	err = agg.forEachBucketElem(agg.buckets[getKeyFn(ts)], func(elem *RunSession) error {
		return errors.New("oops, something went wrong")
	})
	assert.Error(t, err)
}

func TestBucketFuncsFromSize(t *testing.T) {
	ts := time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC)

	// hours
	getKeyFn, nextKeyFn := bucketFuncsFromSize(HourBucketSize)

	key := getKeyFn(ts)
	assert.Equal(t, "2006-01-02T15", key)

	key = nextKeyFn(key)
	assert.Equal(t, "2006-01-02T16", key)

	// days
	getKeyFn, nextKeyFn = bucketFuncsFromSize(DayBucketSize)

	key = getKeyFn(ts)
	assert.Equal(t, "2006-01-02", key)

	key = nextKeyFn(key)
	assert.Equal(t, "2006-01-03", key)

	// weeks
	getKeyFn, nextKeyFn = bucketFuncsFromSize(WeekBucketSize)

	key = getKeyFn(ts)
	assert.Equal(t, "2006-01-02", key)

	key = nextKeyFn(key)
	assert.Equal(t, "2006-01-09", key)

	// months
	getKeyFn, nextKeyFn = bucketFuncsFromSize(MonthBucketSize)

	key = getKeyFn(ts)
	assert.Equal(t, "2006-01", key)

	key = nextKeyFn(key)
	assert.Equal(t, "2006-02", key)
}
