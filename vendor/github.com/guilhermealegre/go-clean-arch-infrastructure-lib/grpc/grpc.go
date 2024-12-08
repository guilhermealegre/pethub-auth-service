package grpc

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	serviceWriter "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/service_writer"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/grpc/interceptor"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	grpcConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/grpc/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Grpc grpc service
type Grpc struct {
	// Name
	name string
	// Configurations
	configs *grpcConfig.Configs
	// App
	app domain.IApp
	// Server
	server *grpc.Server
	// Clients
	clients map[string]*grpc.ClientConn
	// Controllers
	controllers []domain.IController
	// Initialized Server
	initializedServer bool
	// Initialized Clients
	initializedClients bool
	// Writer
	writer io.Writer
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile grpc config file
	configFile  = "grpc.yaml"
	serviceName = "GRPC"
)

// New creates a new grpc service
func New(app domain.IApp, configs *grpcConfig.Configs) *Grpc {
	grpcObj := &Grpc{
		name:    serviceName,
		app:     app,
		clients: map[string]*grpc.ClientConn{},
		writer:  serviceWriter.NewServiceWriter(serviceName),
	}

	if configs != nil {
		grpcObj.configs = configs
	}
	return grpcObj
}

// Name gets the service name
func (g *Grpc) Name() string {
	var text []string
	if g.initializedServer {
		text = append(text, fmt.Sprintf("%s server ready : %d", g.name, g.configs.Server.Port))
	}

	if g.initializedClients {
		for _, clientConfig := range g.configs.Clients {
			text = append(text, fmt.Sprintf("%s client ready [ %s : %d ]", g.name, clientConfig.Name, clientConfig.Port))
		}
	}

	return strings.Join(text, "\n")
}

// Start starts a grpc service
func (g *Grpc) Start() (err error) {
	// initialize configs
	if g.configs == nil {
		if err = g.initConfigs(); err != nil {
			return err
		}
	}

	// initialize server
	err = g.InitServer()
	if err != nil {
		return err
	}

	// initialize clients
	err = g.InitClients()
	if err != nil {
		return err
	}

	g.started = true

	return nil
}

// Stop stops a grpc service
func (g *Grpc) Stop() (err error) {
	if g.initializedServer {
		g.server.GracefulStop()
	}

	if g.initializedClients {
		for _, c := range g.clients {
			if err = c.Close(); err != nil {
				return err
			}
		}
	}

	g.started = false

	return nil
}

// Config gets the grpc configurations
func (g *Grpc) Config() *grpcConfig.Configs {
	return g.configs
}

// ConfigFile gets the configuration file
func (g *Grpc) ConfigFile() string {
	return configFile
}

// initConfigs initialize the configurations
func (g *Grpc) initConfigs() (err error) {
	g.configs = &grpcConfig.Configs{}
	g.configs.AdditionalConfig = g.additionalConfigType
	if err = config.Load(g.ConfigFile(), g.configs); err != nil {
		err = errorCodes.ErrorLoadingConfigFile().Formats(g.ConfigFile(), err)
		message.ErrorMessage(g.Name(), err)
		return err
	}
	return nil
}

// InitServer initialize the grpc server
func (g *Grpc) InitServer() (err error) {
	if g.configs == nil {
		if err = g.initConfigs(); err != nil {
			return err
		}
	}

	if g.initializedServer {
		return nil
	}

	// if no server is needed bypass
	if g.configs.Server.Name == "" {
		return
	}

	// listen the host and port defined in configs
	lis, err := net.Listen("tcp", g.getServerUrl())
	if err != nil {
		return err
	}

	// create server with interceptor
	g.server = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorServerInterceptor(),
			interceptor.PrintServerInterceptor(g.writer),
		),
		grpc.ChainStreamInterceptor(
			interceptor.ErrorServerStreamingInterceptor(),
			interceptor.PrintServerStreamingInterceptor(g.writer),
		),
	)

	// register implementations of the server
	for _, controller := range g.controllers {
		controller.Register()
	}

	g.printServerServices()

	// Register reflection service on gRPC server.
	reflection.Register(g.server)

	// create go routine to listen grpc
	status := make(chan error, 1)
	defer close(status)
	go func(status chan error) {
		if err = g.server.Serve(lis); err != nil {
			status <- err
		}
	}(status)

	g.initializedServer = true

	select {
	case err = <-status:
		return err
	case <-time.After(2 * time.Second):
		return nil
	}
}

func (g *Grpc) printServerServices() {
	for name, desc := range g.server.GetServiceInfo() {
		for _, method := range desc.Methods {
			_, _ = g.writer.Write([]byte(fmt.Sprintf("%s/%s --> %s\n", name, method.Name, desc.Metadata)))
		}
	}
}

// InitClients initialize the clients
func (g *Grpc) InitClients() (err error) {
	if g.configs == nil {
		if err = g.initConfigs(); err != nil {
			return err
		}
	}

	if g.initializedClients {
		return nil
	}

	// if no clients defined , then bypass
	if len(g.configs.Clients) == 0 {
		return
	}

	for _, clientConfig := range g.Config().Clients {
		newClient, err := grpc.NewClient(
			g.getClientUrl(clientConfig.Name),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithChainUnaryInterceptor(
				interceptor.ErrorClientInterceptor(),
				interceptor.PrintClientInterceptor(g.writer),
			),
			grpc.WithChainStreamInterceptor(
				interceptor.ErrorClientStreamingInterceptor(),
				interceptor.PrintClientStreamingInterceptor(g.writer),
			),
		)
		if err != nil {
			return err
		}

		if client, ok := g.clients[clientConfig.Name]; !ok {
			g.clients[clientConfig.Name] = newClient
		} else {
			*client = *newClient //nolint:all
		}
	}

	g.initializedClients = true

	return nil
}

// GetServer gets the server
func (g *Grpc) GetServer() (conn *grpc.Server, err error) {
	return g.server, nil
}

// GetClient gets the client by name
func (g *Grpc) GetClient(name string) *grpc.ClientConn {

	for clientName, c := range g.clients {
		if clientName == name {
			return c
		}
	}

	client := &grpc.ClientConn{}
	g.clients[name] = client

	return client
}

// getServerUrl get url server
func (g *Grpc) getServerUrl() string {
	return g.configs.Server.Host + ":" + strconv.Itoa(g.configs.Server.Port)
}

// getClientUrl gets the url client
func (g *Grpc) getClientUrl(name string) string {

	for _, c := range g.configs.Clients {
		if c.Name == name {
			return c.Host + ":" + strconv.Itoa(c.Port)
		}
	}
	return ""
}

// WithController adds a new controller to the server
func (g *Grpc) WithController(controller domain.IController) domain.IGrpc {
	g.controllers = append(g.controllers, controller)
	return g
}

// WithAdditionalConfigType sets an additional config type
func (g *Grpc) WithAdditionalConfigType(obj interface{}) domain.IGrpc {
	g.additionalConfigType = obj
	return g
}

// Started true if started
func (g *Grpc) Started() bool {
	return g.started
}
