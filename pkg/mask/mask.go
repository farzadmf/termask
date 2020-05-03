package mask

import (
	"io"
)

const (
	valuePattern = `[a-zA-Z0-9%;=/._-]`
)

// Config is used to specify the reader to mask and writer to write the output
type Config struct {
	Reader io.Reader
	Writer io.Writer
}

// NewConfig creates a new mask configuration using the specified reader and writer
func NewConfig(r io.Reader, w io.Writer) Config {
	return Config{
		Reader: r,
		Writer: w,
	}
}

// Masker receives a reader and a writer; reads from the reader and writes masked output to the writer
type Masker interface {
	Mask(config Config)
}
