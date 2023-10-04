package module

import "time"

func CountElapsed(start time.Time) time.Duration {
	return time.Since(start)
}
