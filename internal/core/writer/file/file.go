package file

import (
	"fmt"
	"os"

	"github.com/saromanov/heartbeat/internal/core/writer"
)

// File defines file writer
type File struct {
	file *os.File
}

// New provides initialization of the file writer
func New(fileName string) (writer.Writer, error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	return &File{
		file: f,
	}, nil
}

// Write provides writing of teh data to file
func (f *File) Write(data []byte) error {
	_, err := f.file.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write data: %v", err)
	}
	return nil
}
