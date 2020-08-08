package api

import (
	"time"

	"github.com/saromanov/heartbeat/internal/core"
	"github.com/saromanov/heartbeat/internal/core/writer/stdout"
)

// Heartbeat defines main object
type Heartbeat struct {
	check *core.Check
}

// New provides initialization of the heartbeat
func New() *Heartbeat {
	return &Heartbeat{
		check: core.New(stdout.New()),
	}
}

// AddCheck provides adding of the check
func (h *Heartbeat) AddCheck(title, url string) error {
	return h.check.AddHTTPCheck(core.HTTPCheck{
		Title: title,
		URL:   url,
	})
}

// Run provides starting of the app
func (h *Heartbeat) Run(d time.Duration) {
	h.check.Run(d)
}

// Stats return stats for all registered endpoints
func (h *Heartbeat) Stats() map[int]core.Stats {
	return h.check.Stats()
}
