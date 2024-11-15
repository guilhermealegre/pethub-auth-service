package http

import (
	"bytes"
	"io"
)

type reader struct {
	buffer         *bytes.Buffer
	bytes          []byte
	allowBodyClose bool
}

func NewReader(ioReader io.Reader, allowBodyClose bool) *reader {
	read, _ := io.ReadAll(ioReader)
	return &reader{
		buffer:         bytes.NewBuffer(read),
		bytes:          read,
		allowBodyClose: allowBodyClose,
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.buffer.Read(p)
	if !r.allowBodyClose {
		if (err != nil && err == io.EOF) || n <= len(p) {
			r.buffer.Reset()
			_, _ = r.buffer.Write(r.bytes)
			if n <= len(p) {
				return n, io.EOF
			}
		}
	}
	return n, err
}

func (r *reader) Close() error {
	if r.allowBodyClose {
		r.buffer.Reset()
	}
	return nil
}
