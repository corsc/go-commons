package iocloser

import (
	"io"
)

// Logger defines the logging function for this package
type Logger func(template string, args ...interface{})

// Close will close the supplied ReadCloser and optionally log
func Close(reader io.Closer, logger ...Logger) {
	err := reader.Close()
	if err != nil {
		if len(logger) > 0 {
			logger[0]("error while closing reader. err: %s", err)
		}
	}
}
