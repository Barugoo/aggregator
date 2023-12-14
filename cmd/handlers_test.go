package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Barugoo/twaiv/internal/aggregator"
	"github.com/Barugoo/twaiv/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAggregateWorkoutsHandle(t *testing.T) {
	testcases := []struct {
		name                 string
		nweeks               int
		requestBody          []models.RunSession
		expectedStatus       int
		expectedResponseBody models.WorkoutAnalyseResponse
	}{
		{
			name:           "normal",
			nweeks:         2,
			expectedStatus: http.StatusOK,
			requestBody: []models.RunSession{
				{
					Distance:  100,
					Time:      models.DurationSec{Duration: time.Second * 1000},
					Timestamp: models.TimeRFC3339{Time: time.Now()},
				},
				{
					Distance:  150,
					Time:      models.DurationSec{Duration: time.Second * 2000},
					Timestamp: models.TimeRFC3339{Time: time.Now().AddDate(0, 0, 0)},
				},
				{
					Distance:  200,
					Time:      models.DurationSec{Duration: time.Second * 3000},
					Timestamp: models.TimeRFC3339{Time: time.Now().AddDate(0, 0, 0)},
				},
				{
					Distance:  250,
					Time:      models.DurationSec{Duration: time.Second * 4000},
					Timestamp: models.TimeRFC3339{Time: time.Now().AddDate(0, -4, 0)},
				},
			},
			expectedResponseBody: models.WorkoutAnalyseResponse{
				MediumDistance:       150,
				MediumTime:           models.DurationSec{Duration: time.Second * 2000},
				MediumWeeklyDistance: 450,
				MediumWeeklyTime:     models.DurationSec{Duration: time.Second * 6000},
				MaxDistance:          200,
				MaxTime:              models.DurationSec{Duration: time.Second * 3000},
				MaxWeeklyDistance:    450,
				MaxWeeklyTime:        models.DurationSec{Duration: time.Second * 6000},
			},
		},
	}

	for _, tc := range testcases {
		inputCh := make(chan *aggregator.RunSession)

		agg, err := aggregator.NewInMemory(inputCh, aggregator.WeekBucketSize)
		assert.NoError(t, err)

		h := AggregateWorkoutsHandle(inputCh, agg)

		srv := httptest.NewServer(h)

		byt, err := json.Marshal(tc.requestBody)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, srv.URL+"?nweeks="+strconv.Itoa(tc.nweeks), bytes.NewReader(byt))
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedStatus, resp.StatusCode)

		var actualBody models.WorkoutAnalyseResponse
		err = json.NewDecoder(resp.Body).Decode(&actualBody)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedResponseBody, actualBody)
	}
}
