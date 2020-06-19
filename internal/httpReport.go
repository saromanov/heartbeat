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
	Error  string `json:"error"`
	Url    string `json:"url"`
	Tags   []string `json:"tags"`
}