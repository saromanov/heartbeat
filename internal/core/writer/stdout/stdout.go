package stdout

import (
	"fmt"

	"github.com/saromanov/heartbeat/internal/core/writer"
)

// Stdout provides writing data to output
type Stdout struct {
}

// New in that case its puppet initialization
func New() writer.Writer {
	return &Stdout{}
}

// Write defines output for data
func (s *Stdout) Write(data []byte) error {
	fmt.Println(string(data))
	return nil
}
