package writer

// Writer provides writing of the data to sources
type Writer interface {
	Write([]byte) error
}
