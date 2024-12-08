package service_writer

import (
	"fmt"
	"io"

	loggerWriter "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/writer"
)

const (
	debugMessage = "[%s-debug]"
	separator    = " "
)

type ServiceWriter struct {
	name string
	io.Writer
}

func NewServiceWriter(name string) *ServiceWriter {
	return &ServiceWriter{
		name:   name,
		Writer: loggerWriter.NewConsole(),
	}
}

func (w *ServiceWriter) Write(p []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(fmt.Sprint(fmt.Sprintf(debugMessage, w.name), separator)), p...))
}
