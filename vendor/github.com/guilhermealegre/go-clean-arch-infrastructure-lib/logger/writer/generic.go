package writer

import (
	"io"
)

// Generic
type Generic struct {
	// handler
	handler WriterHandler
	// fallback
	fallback io.Writer
}

// NewGeneric creates a new writer to generic
func NewGeneric(handler WriterHandler, fallback io.Writer) *Generic {
	generic := &Generic{
		handler:  handler,
		fallback: fallback,
	}

	return generic
}

// Write writes the bytes to the writer
func (r *Generic) Write(message []byte) (n int, err error) {
	if r.handler != nil {
		// execute handler
		err = r.handler(message)
	}

	if (r.handler == nil || err != nil) && r.fallback != nil {
		if n, err = r.fallback.Write(message); err != nil {
			return 0, err
		}
		return n, nil
	}

	return len(message), nil
}
