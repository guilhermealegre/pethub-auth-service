package v1

import (
	"github.com/stretchr/testify/mock"
)

func NewStreamingMock() *StreamingMock {
	return &StreamingMock{}
}

type StreamingMock struct {
	mock.Mock
}
