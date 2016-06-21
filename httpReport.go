package heartbeat

import (
	"time"
)
type HTTPReport struct {
	Items []HTTPItem `json:"items"`
}

type HTTPItem struct {
	Name  string `json:"name"`
	Status  string `json:"status"`
	StatusCode string `json:"statusCode"`
	Time   *time.Time `json:"time"`
	Tags   []string `json:"tags"`
}