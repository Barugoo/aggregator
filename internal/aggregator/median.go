package aggregator

import (
	"container/list"
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func (a *aggregator[T]) Median(from, to time.Time) (*AggregationResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var res AggregationResult

	var distances []int64
	var times []time.Duration

	var bucketTotalDistances []int64
	var bucketTotalTimes []time.Duration
	err := a.forEachBucket(from, to, func(bucketElems *list.List) error {

		startBucketIdx := len(distances)
		err := a.forEachBucketElem(bucketElems, func(elem T) error {

			distances = append(distances, elem.GetDistance())
			times = append(times, elem.GetDuration())
			return nil
		})
		if err != nil {
			return fmt.Errorf("unable to calculate median: %w", err)
		}

		bucketTotalDistances = append(bucketTotalDistances, sum(distances[startBucketIdx:]))
		bucketTotalTimes = append(bucketTotalTimes, sum(times[startBucketIdx:]))

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to calculate median: %w", err)
	}

	res.Distance = int64(median(distances))
	res.Duration = time.Duration(median(times))
	res.DistanceByBucket = int64(median(bucketTotalDistances))
	res.DurationByBucket = time.Duration(median(bucketTotalTimes))

	return &res, nil
}

// to save time I copied this function body from here: https://gosamples.dev/calculate-median/
func median[T constraints.Float | constraints.Integer](data []T) float64 {
	dataCopy := make([]T, len(data))
	copy(dataCopy, data)

	slices.Sort(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = float64((dataCopy[l/2-1] + dataCopy[l/2]) / 2.0)
	} else {
		median = float64(dataCopy[l/2])
	}

	return median
}

func sum[T constraints.Float | constraints.Integer](data []T) (res T) {
	for _, elem := range data {
		res += elem
	}
	return res
}
