package app

import (
	"fmt"
	instanceStateMachine "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/state_machine/instance"
	"os"
	"os/signal"
	"syscall"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"

	appConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/app/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

// App keeps all the tools that the lib handles
type App struct {
	// Configuration
	config *appConfig.Config

	// Logger
	logger domain.ILogger
	// Database
	database domain.IDatabase
	// Rabbitmq
	rabbitmq domain.IRabbitMQ
	// SQS
	sqs domain.ISQS
	// Redis
	redis domain.IRedis
	// ElasticSearch
	elasticSearch domain.IElasticSearch
	// Grpc
	grpc domain.IGrpc
	// Http
	http domain.IHttp
	// Validator
	validator domain.IValidator
	// Aws
	aws domain.IAws
	// S3
	s3 domain.IS3
	// Tracer
	tracer domain.ITracer
	// Meter
	meter domain.IMeter
	// Datatable
	datatable domain.IDatatable
	// State Machine
	stateMachine instanceStateMachine.IStateMachineService
	// Services
	services []domain.IService
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	configFile = "app.yaml"
)

// New creates a new App instance
func New(config *appConfig.Config) domain.IApp {
	app := &App{}

	if config != nil {
		app.config = config
	}

	return app
}

// Name app name
func (a *App) Name() string {
	return a.config.Name
}

// Start starts the app lib
func (a *App) Start() (err error) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if a.config == nil {
		a.config = &appConfig.Config{}
		a.config.AdditionalConfig = a.additionalConfigType
		if err = config.Load(a.ConfigFile(), a.config); err != nil {
			err = errorCodes.ErrorLoadingConfigFile().Formats(a.ConfigFile(), err)
			message.ErrorMessage(a.config.Name, err)
			return err
		}
	}

	for _, s := range a.services {
		if err = s.Start(); err != nil {
			message.ErrorMessage(s.Name(), err)
			_ = a.Stop() // stop previous started services
			return err
		} else {
			message.StartMessage(s.Name())
		}
	}

	a.started = true
	message.StartMessage(fmt.Sprintf("Service [%s]", a.config.Name))

	// wait for a TERM signal
	<-exit

	return a.Stop()
}

// Stop stops the app services
func (a *App) Stop() (err error) {
	var errStop error

	for i := len(a.services) - 1; i >= 0; i-- {
		if a.services[i].Started() {
			message.StopMessage(a.services[i].Name())
			if err = a.services[i].Stop(); err != nil {
				// stop every service ignoring the returned errors
				errStop = err
				message.ErrorMessage(a.services[i].Name(), err)
			}
		}
	}

	a.started = false
	message.StopMessage(fmt.Sprintf("Service [%s]", a.config.Name))

	return errStop
}

// WithLogger sets the logger service
func (a *App) WithLogger(logger domain.ILogger) domain.IApp {
	a.addService(logger)
	a.logger = logger
	return a
}

// Logger gets the logger service
func (a *App) Logger() domain.ILogger {
	return a.logger
}

// WithDatabase sets the database service
func (a *App) WithDatabase(database domain.IDatabase) domain.IApp {
	a.addService(database)
	a.database = database
	return a
}

// Database gets the database service
func (a *App) Database() domain.IDatabase {
	return a.database
}

// WithRabbitmq sets the rabbitmq service
func (a *App) WithRabbitmq(rabbitmq domain.IRabbitMQ) domain.IApp {
	a.addService(rabbitmq)
	a.rabbitmq = rabbitmq
	return a
}

// Rabbitmq gets the rabbitmq service
func (a *App) Rabbitmq() domain.IRabbitMQ {
	return a.rabbitmq
}

// WithSQS sets the sqs service
func (a *App) WithSQS(sqs domain.ISQS) domain.IApp {
	a.addService(sqs)
	a.sqs = sqs
	return a
}

// SQS gets the sqs service
func (a *App) SQS() domain.ISQS {
	return a.sqs
}

// WithRedis sets the redis service
func (a *App) WithRedis(redis domain.IRedis) domain.IApp {
	a.addService(redis)
	a.redis = redis
	return a
}

// Redis gets the redis service
func (a *App) Redis() domain.IRedis {
	return a.redis
}

// WithElasticSearch sets the elastic search service
func (a *App) WithElasticSearch(elasticSearch domain.IElasticSearch) domain.IApp {
	a.addService(elasticSearch)
	a.elasticSearch = elasticSearch
	return a
}

// ElasticSearch gets the elastic search service
func (a *App) ElasticSearch() domain.IElasticSearch {
	return a.elasticSearch
}

// WithGrpc sets the grpc service
func (a *App) WithGrpc(grpc domain.IGrpc) domain.IApp {
	a.addService(grpc)
	a.grpc = grpc
	return a
}

// Grpc gets the grpc service
func (a *App) Grpc() domain.IGrpc {
	return a.grpc
}

// WithHttp sets the http service
func (a *App) WithHttp(http domain.IHttp) domain.IApp {
	a.addService(http)
	a.http = http
	return a
}

// Http gets the http service
func (a *App) Http() domain.IHttp {
	return a.http
}

// WithValidator sets the validator service
func (a *App) WithValidator(validator domain.IValidator) domain.IApp {
	a.addService(validator)
	a.validator = validator
	return a
}

func (a *App) Validator() domain.IValidator {
	return a.validator
}

// WithAws sets the aws service
func (a *App) WithAws(aws domain.IAws) domain.IApp {
	a.addService(aws)
	a.aws = aws
	return a
}

// Aws gets the aws service
func (a *App) Aws() domain.IAws {
	return a.aws
}

// WithS3 sets the s3 service
func (a *App) WithS3(s3 domain.IS3) domain.IApp {
	a.addService(s3)
	a.s3 = s3
	return a
}

// S3 gets the s3 service
func (a *App) S3() domain.IS3 {
	return a.s3
}

// WithTracer sets the tracer service
func (a *App) WithTracer(tracer domain.ITracer) domain.IApp {
	a.addService(tracer)
	a.tracer = tracer
	return a
}

// Tracer gets the tracer service
func (a *App) Tracer() domain.ITracer {
	return a.tracer
}

// WithMeter sets the meter service
func (a *App) WithMeter(meter domain.IMeter) domain.IApp {
	a.addService(meter)
	a.meter = meter
	return a
}

// Tracer gets the tracer service
func (a *App) Meter() domain.IMeter {
	return a.meter
}

// WithDatatable sets the datatable
func (a *App) WithDatatable(datatable domain.IDatatable) domain.IApp {
	a.addService(datatable)
	a.datatable = datatable
	return a
}

// Datatable gets the datatable
func (a *App) Datatable() domain.IDatatable {
	return a.datatable
}

// WithDatatable sets the datatable
func (a *App) WithStateMachine(stateMachine instanceStateMachine.IStateMachineService) domain.IApp {
	a.addService(stateMachine)
	a.stateMachine = stateMachine
	return a
}

// StateMachine gets the state machine service
func (a *App) StateMachine() instanceStateMachine.IStateMachineService {
	return a.stateMachine
}

// Config gets the app configurations
func (a *App) Config() *appConfig.Config {
	return a.config
}

// addService adds a service to the app
func (a *App) addService(service domain.IService) {
	if service != nil {
		a.services = append(a.services, service)
	}
}

// ConfigFile gets the configuration file name
func (a *App) ConfigFile() string {
	return configFile
}

// WithAdditionalConfigType sets an additional config type
func (a *App) WithAdditionalConfigType(obj interface{}) domain.IApp {
	a.additionalConfigType = obj
	return a
}

// Started true if started
func (a *App) Started() bool {
	return a.started
}
