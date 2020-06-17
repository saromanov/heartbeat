package main

import (
	"time"

	"github.com/saromanov/heartbeat"
)

func main() {
	h := heartbeat.New()
	h.AddHTTPCheck(heartbeat.HTTPCheck{Title: "some", URL: "https://github.com"})
	h.Run(1 * time.Second)
}
