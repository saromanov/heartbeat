package main

import (
	"flag"
	"fmt"

	"github.com/saromanov/heartbeat/internal/config"
	"github.com/saromanov/heartbeat/internal/server"
)

func loadConfig(path string) (*config.Config, error) {
	cfg, err := config.Unmarshal(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %v", err)
	}
	return cfg, err
}
func main() {
	path := flag.String("config-path", "", "path to config path")
	flag.Parse()

	conf := config.Default()
	if *path != "" {
		confTmp, err := loadConfig(*path)
		if err != nil {
			panic(err)
		}
		conf = confTmp
	}
	server.Run(conf)
}
