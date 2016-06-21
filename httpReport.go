package heartbeat

import (
	"time"
)
type HTTPReport struct {

}

type HTTPItem struct {
	Name  string
	Status  string
	Time   *time.Time
	Tags   []string
}