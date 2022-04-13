package util

import "time"

type Elapsed time.Duration

func Timer(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start).Truncate(time.Millisecond) / time.Millisecond
}
