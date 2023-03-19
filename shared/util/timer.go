package util

import "time"

type Elapsed time.Duration

func Timer(fn func()) int64 {
	start := time.Now()
	fn()
	return int64(time.Since(start).Truncate(time.Millisecond) / time.Millisecond)
}
