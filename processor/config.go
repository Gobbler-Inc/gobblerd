package processor

import "time"

var (
	taskInterval time.Duration = 1 * time.Second
)

func TaskInterval() time.Duration {
	return taskInterval
}

func SetTaskInterval(newInterval time.Duration) {
	taskInterval = newInterval
}
