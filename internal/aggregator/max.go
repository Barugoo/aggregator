package aggregator

import (
	"github.com/zyedidia/generic/list"

	"fmt"
	"time"
)

func (a *aggregator[T]) Max(from, to time.Time) (*AggregationResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var res AggregationResult
	err := a.forEachBucket(from, to, func(bucketElems *list.List[T]) error {

		var bucketTotalDistance int64
		var bucketTotalTime time.Duration

		bucketElems.Front.Each(func(elem T) {
			if elemTime := elem.GetDuration(); elemTime > res.Duration {
				res.Duration = elemTime
			}
			if elemDistance := elem.GetDistance(); elemDistance > res.Distance {
				res.Distance = elemDistance
			}
			bucketTotalDistance += elem.GetDistance()
			bucketTotalTime += elem.GetDuration()
		})

		if bucketTotalDistance > res.DistanceByBucket {
			res.DistanceByBucket = bucketTotalDistance
		}
		if bucketTotalTime > res.DurationByBucket {
			res.DurationByBucket = bucketTotalTime
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to calculate median: %w", err)
	}

	return &res, nil
}
