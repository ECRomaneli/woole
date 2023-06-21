package timer

import "time"

type Elapsed time.Duration

func Exec(fn func()) int64 {
	start := time.Now()
	fn()
	return int64(time.Since(start))
}
