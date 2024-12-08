package grpc

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	grpcConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/grpc/config"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func NewGrpcMock() *GrpcMock {
	return &GrpcMock{}
}

type GrpcMock struct {
	mock.Mock
}

func (g *GrpcMock) Name() string {
	args := g.Called()
	return args.Get(0).(string)
}

func (g *GrpcMock) Start() error {
	args := g.Called()
	return args.Error(0)
}

func (g *GrpcMock) Stop() error {
	args := g.Called()
	return args.Error(0)
}

func (g *GrpcMock) Config() *grpcConfig.Configs {
	args := g.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*grpcConfig.Configs)
}

func (g *GrpcMock) ConfigFile() string {
	args := g.Called()
	return args.Get(0).(string)
}

func (g *GrpcMock) GetClient(name string) *grpc.ClientConn {
	args := g.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*grpc.ClientConn)
}

func (g *GrpcMock) GetServer() (*grpc.Server, error) {
	args := g.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grpc.Server), args.Error(1)
}

func (g *GrpcMock) WithController(controller domain.IController) domain.IGrpc {
	args := g.Called(controller)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IGrpc)
}

// WithAdditionalConfigType sets an additional config type
func (g *GrpcMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := g.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (g *GrpcMock) Started() bool {
	args := g.Called()
	return args.Get(0).(bool)
}
