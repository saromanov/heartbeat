package main

import (
	"github.com/saromanov/heartbeat/internal/config"
	"github.com/saromanov/heartbeat/internal/server"
)

func main() {
	server.Run(config.Default())
}
