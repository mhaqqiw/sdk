package qentity

import (
	"time"

	"github.com/google/uuid"
)

type UUIDData struct {
	UUID          uuid.UUID
	Timestamp     time.Time
	ClockSequence int
	Domain        string
	NodeID        string
	ID            uint32
	MarshalText   string
	String        string
	Version       string
	URN           string
	Value         string
	Variant       string
}
