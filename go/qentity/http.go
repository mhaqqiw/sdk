package qentity

import (
	"time"
)

type Response struct {
	Status      string        `json:"status"`
	Code        int           `json:"code"`
	Message     interface{}   `json:"data"`
	ProcessTime time.Duration `json:"process_time"`
	Version     string        `json:"version"`
}
