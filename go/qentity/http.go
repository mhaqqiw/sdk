package qentity

import (
	"encoding/json"
	"time"
)

type Response struct {
	Status      string          `json:"status"`
	Code        int             `json:"code"`
	Message     json.RawMessage `json:"data"`
	ProcessTime time.Duration   `json:"process_time"`
	Version     string          `json:"version"`
}
