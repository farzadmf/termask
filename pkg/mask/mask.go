package mask

import (
	"io"
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
