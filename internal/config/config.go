package config

// Config defines configuration for the app
type Config struct {
	Address string  `yaml:"address"`
	Checks  []Check `yaml:"checks"`
}

// Check defines http check
type Check struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
