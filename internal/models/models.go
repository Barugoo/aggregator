package models

import (
	"fmt"
	"strconv"
	"time"
)

type RunSession struct {
	// run distance in meters
	Distance int64 `json:"distance"`
	// run duration
	Time DurationSec `json:"time"`
	// run start ts
	Timestamp TimeRFC3339 `json:"timestamp"`
}

type DurationSec struct {
	time.Duration
}

func (d *DurationSec) UnmarshalJSON(data []byte) error {
	parsedDur, err := time.ParseDuration(string(data) + "s")
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}
	d.Duration = parsedDur
	return nil
}

func (d *DurationSec) MarshalJSON() ([]byte, error) {
	secs := int64(d.Seconds())
	return []byte(strconv.FormatInt(secs, 10)), nil
}

type TimeRFC3339 struct {
	time.Time
}

func (t *TimeRFC3339) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(time.RFC3339, string(data[1:len(data)-1]))
	if err != nil {
		return fmt.Errorf("unable to parse RFC3339 time: %w", err)
	}
	t.Time = parsedTime
	return nil
}

func (t *TimeRFC3339) MarshalJSON() (res []byte, err error) {
	str := "\"" + t.Format(time.RFC3339) + "\""
	return []byte(str), nil
}

type WorkoutAnalyseResponse struct {
	MediumDistance       int64       `json:"medium_distance"`
	MediumTime           DurationSec `json:"medium_time"`
	MaxDistance          int64       `json:"max_distance"`
	MaxTime              DurationSec `json:"max_time"`
	MediumWeeklyDistance int64       `json:"medium_weekly_distance"`
	MediumWeeklyTime     DurationSec `json:"medium_weekly_time"`
	MaxWeeklyDistance    int64       `json:"max_weekly_distance"`
	MaxWeeklyTime        DurationSec `json:"max_weekly_time"`
}
