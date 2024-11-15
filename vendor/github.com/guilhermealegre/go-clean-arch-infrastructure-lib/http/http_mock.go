package http

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	httpConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/http/config"
	"github.com/stretchr/testify/mock"
)

func NewHttpMock() *HttpMock {
	return &HttpMock{}
}

type HttpMock struct {
	mock.Mock
}

func (h *HttpMock) Name() string {
	args := h.Called()
	return args.Get(0).(string)
}

func (h *HttpMock) Start() error {
	args := h.Called()
	return args.Error(0)
}

func (h *HttpMock) Stop() error {
	args := h.Called()
	return args.Error(0)
}

func (h *HttpMock) ConfigFile() string {
	args := h.Called()
	return args.Get(0).(string)
}

func (h *HttpMock) Config() *httpConfig.Config {
	args := h.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*httpConfig.Config)
}

func (h *HttpMock) WithMiddleware(controller domain.IMiddleware) domain.IHttp {
	args := h.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IHttp)
}

func (h *HttpMock) WithController(controller domain.IController) domain.IHttp {
	args := h.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IHttp)
}

func (h *HttpMock) WithRouter(router *gin.Engine) domain.IHttp {
	args := h.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IHttp)
}

func (h *HttpMock) Router() *gin.Engine {
	args := h.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gin.Engine)
}

// WithAdditionalConfigType sets an additional config type
func (h *HttpMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := h.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (h *HttpMock) Started() bool {
	args := h.Called()
	return args.Get(0).(bool)
}
