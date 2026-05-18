package observability

import "time"

type Timing struct {
	Started time.Time
	Ended   time.Time
}

func (t Timing) Duration() time.Duration {
	if t.Ended.Before(t.Started) {
		return 0
	}
	return t.Ended.Sub(t.Started)
}
