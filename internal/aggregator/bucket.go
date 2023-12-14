package aggregator

import (
	"container/list"
	"fmt"
	"time"
)

type BucketSize uint8

const (
	HourBucketSize BucketSize = iota
	DayBucketSize
	WeekBucketSize
	MonthBucketSize
)

type bucketKeyFn func(time.Time) string
type nextKeyFn func(string) string

func (a *aggregator[T]) forEachBucket(from, to time.Time, fn func(bucketElems *list.List) error) error {
	toKey := a.bucketKeyFn(to)

	for fromKey := a.bucketKeyFn(from); fromKey <= toKey; fromKey = a.nextKeyFn(fromKey) {
		bucketElems, ok := a.buckets[fromKey]
		if !ok {
			continue
		}
		if err := fn(bucketElems); err != nil {
			return fmt.Errorf("forEachBucket: fn returned error: %w", err)
		}
	}
	return nil
}

func (a *aggregator[T]) forEachBucketElem(bucketElems *list.List, fn func(elem T) error) error {
	for bucketElem := bucketElems.Front(); bucketElem != nil; bucketElem = bucketElem.Next() {
		elem, ok := bucketElem.Value.(T)
		if !ok {
			return fmt.Errorf("forEachBucketElem: the bucket list should store object of %T type", elem)
		}
		if err := fn(elem); err != nil {
			return fmt.Errorf("forEachBucketElem: fn returned error: %w", err)
		}
	}
	return nil
}

func bucketFuncsFromSize(size BucketSize) (getKey bucketKeyFn, getNextKey nextKeyFn) {
	switch size {
	case HourBucketSize:
		keyFormat := "2006-01-02T15"

		getKey = func(t time.Time) string {
			return t.Format(keyFormat)
		}
		getNextKey = func(s string) string {
			t, _ := time.Parse(keyFormat, s)
			// we can omit err check since we create these keys
			return t.Add(time.Hour).Format(keyFormat)
		}

	case DayBucketSize:
		keyFormat := "2006-01-02"

		getKey = func(t time.Time) string {
			return t.Format(keyFormat)
		}
		getNextKey = func(s string) string {
			t, _ := time.Parse(keyFormat, s)
			// we can omit err check since we create these keys
			return t.AddDate(0, 0, 1).Format(keyFormat)
		}

	case WeekBucketSize:
		keyFormat := "2006-01-02"

		getKey = func(t time.Time) string {
			newT := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			newT = newT.AddDate(0, 0, int(time.Monday-newT.Weekday()))
			return newT.Format(keyFormat)
		}
		getNextKey = func(s string) string {
			t, _ := time.Parse(keyFormat, s)
			// we can omit err check since we create these keys
			return t.AddDate(0, 0, 7).Format(keyFormat)
		}

	case MonthBucketSize:
		keyFormat := "2006-01"

		getKey = func(t time.Time) string {
			return t.Format(keyFormat)
		}

		getNextKey = func(s string) string {
			t, _ := time.Parse(keyFormat, s)
			// we can omit err check since we create these keys
			return t.AddDate(0, 1, 0).Format(keyFormat)
		}
	}
	return getKey, getNextKey
}
