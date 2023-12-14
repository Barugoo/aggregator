package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Barugoo/twaiv/internal/aggregator"
	"github.com/Barugoo/twaiv/internal/models"
)

func AggregateWorkoutsHandle(inputCh chan<- *aggregator.RunSession, agg aggregator.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weeksNumberStr := r.URL.Query().Get("nweeks")
		if weeksNumberStr == "" {
			http.Error(w, "nweeks query-param should be specified", http.StatusBadRequest)
			return
		}

		weeksNumber, err := strconv.Atoi(weeksNumberStr)
		if err != nil {
			http.Error(w, "nweeks query-param should be an integer number", http.StatusBadRequest)
			return
		}

		var sess []*models.RunSession
		if err := json.NewDecoder(r.Body).Decode(&sess); err != nil {
			http.Error(w, "unable to deserialize json body: "+err.Error(), http.StatusBadRequest)
			return
		}

		for _, ses := range sess {
			aggSess := aggregator.RunSession{
				Timestamp: ses.Timestamp.Time,
				Duration:  ses.Time.Duration,
				Distance:  ses.Distance,
			}
			// made this asynchronous for two reasons
			// 1. to show-off
			// 2. to make input API as flexible as possible, being pretty easy to pipe it with some other async data streams (mq, db polling etc.)
			inputCh <- &aggSess
		}

		// since we have RWMutex insde inmemory aggregator implementation and inputCh has no buffer
		// it is guaranteed that by the time we reach this line at least all the request body RunSessions are consumed
		resMax, err := agg.Max(time.Now().AddDate(0, 0, -weeksNumber*7), time.Now())
		if err != nil {
			log.Printf("agg.Max returned error: %v", err)
			http.Error(w, "unable to aggregate provided run sessions", http.StatusBadRequest)
			return
		}

		resMedian, err := agg.Median(time.Now().AddDate(0, 0, -weeksNumber*7), time.Now())
		if err != nil {
			log.Printf("agg.Median returned error: %v", err)
			http.Error(w, "unable to aggregate provided run sessions", http.StatusBadRequest)
			return
		}

		resp := models.WorkoutAnalyseResponse{
			MaxDistance:       resMax.Distance,
			MaxTime:           models.DurationSec{Duration: resMax.Duration},
			MaxWeeklyDistance: resMax.DistanceByBucket,
			MaxWeeklyTime:     models.DurationSec{Duration: resMax.DurationByBucket},

			MediumDistance:       resMedian.Distance,
			MediumTime:           models.DurationSec{Duration: resMedian.Duration},
			MediumWeeklyDistance: resMedian.DistanceByBucket,
			MediumWeeklyTime:     models.DurationSec{Duration: resMedian.DurationByBucket},
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

	}
}
