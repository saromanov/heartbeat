package config

import (
	"time"

	"github.com/saromanov/cowrow"
)

// Config defines configuration for the app
type Config struct {
	Duration time.Duration `yaml:"duration"`
	Address  string        `yaml:"address"`
	Checks   []Check       `yaml:"checks"`
}

// Check defines http check
type Check struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Unmarshal provides unmarshaling of the config
func Unmarshal(path string) (*Config, error) {
	var c *Config
	if err := cowrow.LoadByPath(path, &c); err != nil {
		return nil, err
	}
	return c, nil
}
