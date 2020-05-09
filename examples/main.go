package main

import (
	"fmt"

	"github.com/saromanov/heartbeat"
)

func main() {
	h := heartbeat.New()
	h.AddHTTPCheck("some", "https://github.com")
	r, err := h.CheckHTTP()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
