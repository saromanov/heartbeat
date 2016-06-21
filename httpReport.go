package heartbeat

import (
	"time"
)
type HTTPReport struct {
	items []HTTPItem `json:"items"`
}

type HTTPItem struct {
	Name  string `json:"name"`
	Status  string `json:"status"`
	Time   *time.Time `json:"time"`
	Tags   []string `json:"tags"`
}