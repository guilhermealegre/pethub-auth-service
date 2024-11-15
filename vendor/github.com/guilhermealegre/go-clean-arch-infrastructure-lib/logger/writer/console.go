package writer

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

// File
type Console struct {
	writer io.Writer
}

// NewConsole creates a new writer to Console
func NewConsole() *Console {
	return &Console{
		writer: zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
	}
}

// Write write bytes to file
func (f *Console) Write(bytes []byte) (n int, err error) {
	if n, err = os.Stdout.Write(bytes); err != nil {
		err = errorCodes.ErrorWritingConsole()
		fmt.Println(err)
		return 0, err
	}

	return n, nil
}
