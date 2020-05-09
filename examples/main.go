package main

import (
	"time"

	"github.com/saromanov/heartbeat"
)

func main() {
	h := heartbeat.New()
	h.AddHTTPCheck("some", "https://github.com")
	h.Run(1 * time.Second)
}
