package writer

import "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"

type Fallback struct {
	reader domain.FallbackReader
	writer domain.FallbackWriter
}

type WriterHandler func(message []byte) error
