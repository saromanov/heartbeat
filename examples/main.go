package main

import (
	"time"

	"github.com/saromanov/heartbeat"
)

func main() {
	h := heartbeat.New()
	h.AddHTTPCheck(heartbeat.HTTPCheck{Title: "some", URL: "https://github.com"})
	h.AddHTTPCheck(heartbeat.HTTPCheck{Title: "ya", URL: "https://ya.ru"})
	h.Run(1 * time.Second)
}
