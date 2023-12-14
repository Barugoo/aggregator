package aggregator

import (
	"time"
)

type Row interface {
	GetDistance() int64
	GetDuration() time.Duration
	GetTimestamp() time.Time
}

type RunSession struct {
	Distance  int64
	Duration  time.Duration
	Timestamp time.Time
}

func (ses *RunSession) GetDuration() time.Duration {
	return ses.Duration
}

func (ses *RunSession) GetDistance() int64 {
	return ses.Distance
}

func (ses *RunSession) GetTimestamp() time.Time {
	return ses.Timestamp
}
