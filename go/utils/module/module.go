package module

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CountElapsed(start time.Time) time.Duration {
	return time.Since(start)
}

func GetDataFromUUIDV1(uuid uuid.UUID) (time.Time, int, string) {
	t := uuid.Time()
	sec, nsec := t.UnixTime()
	timeStamp := time.Unix(sec, nsec)

	// Extract the clock sequence (14 bits)
	clockSequence := uuid.ClockSequence()

	// Extract the node (48 bits)
	node := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", uuid[10], uuid[11], uuid[12], uuid[13], uuid[14], uuid[15])

	return timeStamp, clockSequence, node
}
