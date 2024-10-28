package app

import (
	appConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/app/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	instanceStateMachine "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/state_machine/instance"
	"github.com/stretchr/testify/mock"
)

func NewAppMock() *AppMock {
	return &AppMock{}
}

type AppMock struct {
	mock.Mock
}

func (a *AppMock) Name() string {
	args := a.Called()
	return args.Get(0).(string)
}

// Start starts the app lib
func (a *AppMock) Start() error {
	args := a.Called()
	return args.Error(0)
}

// Stop stops the app services
func (a *AppMock) Stop() error {
	args := a.Called()
	return args.Error(0)
}

// WithLogger sets the logger service
func (a *AppMock) WithLogger(logger domain.ILogger) domain.IApp {
	args := a.Called(logger)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Logger gets the logger service
func (a *AppMock) Logger() domain.ILogger {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ILogger)
}

// WithDatabase sets the database service
func (a *AppMock) WithDatabase(database domain.IDatabase) domain.IApp {
	args := a.Called(database)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Database gets the database service
func (a *AppMock) Database() domain.IDatabase {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IDatabase)
}

// WithRabbitmq sets the rabbitmq service
func (a *AppMock) WithRabbitmq(rabbitmq domain.IRabbitMQ) domain.IApp {
	args := a.Called(rabbitmq)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Rabbitmq gets the rabbitmq service
func (a *AppMock) Rabbitmq() domain.IRabbitMQ {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IRabbitMQ)
}

// WithSQS sets the sqs service
func (a *AppMock) WithSQS(sqs domain.ISQS) domain.IApp {
	args := a.Called(sqs)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// SQS gets the sqs service
func (a *AppMock) SQS() domain.ISQS {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ISQS)
}

// WithRedis sets the redis service
func (a *AppMock) WithRedis(redis domain.IRedis) domain.IApp {
	args := a.Called(redis)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Redis gets the redis service
func (a *AppMock) Redis() domain.IRedis {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IRedis)
}

// WithElasticSearch sets the elastic search service
func (a *AppMock) WithElasticSearch(elasticSearch domain.IElasticSearch) domain.IApp {
	args := a.Called(elasticSearch)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// ElasticSearch gets the elastic search service
func (a *AppMock) ElasticSearch() domain.IElasticSearch {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IElasticSearch)
}

// WithGrpc sets the grpc service
func (a *AppMock) WithGrpc(grpc domain.IGrpc) domain.IApp {
	args := a.Called(grpc)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Grpc gets the grpc service
func (a *AppMock) Grpc() domain.IGrpc {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IGrpc)
}

// WithHttp sets the http service
func (a *AppMock) WithHttp(http domain.IHttp) domain.IApp {
	args := a.Called(http)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Http gets the http service
func (a *AppMock) Http() domain.IHttp {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IHttp)
}

// WithValidator sets the validator service
func (a *AppMock) WithValidator(validator domain.IValidator) domain.IApp {
	args := a.Called(validator)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

func (a *AppMock) Validator() domain.IValidator {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IValidator)
}

// WithAws sets the aws service
func (a *AppMock) WithAws(aws domain.IAws) domain.IApp {
	args := a.Called(aws)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Aws gets the aws service
func (a *AppMock) Aws() domain.IAws {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IAws)
}

// WithS3 sets the s3 service
func (a *AppMock) WithS3(s3 domain.IS3) domain.IApp {
	args := a.Called(s3)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// S3 gets the s3 service
func (a *AppMock) S3() domain.IS3 {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IS3)
}

// WithTracer sets the tracer service
func (a *AppMock) WithTracer(tracer domain.ITracer) domain.IApp {
	args := a.Called(tracer)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Tracer gets the tracer service
func (a *AppMock) Tracer() domain.ITracer {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ITracer)
}

// WithMeter sets the meter service
func (a *AppMock) WithMeter(meter domain.IMeter) domain.IApp {
	args := a.Called(meter)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Meter gets the meter service
func (a *AppMock) Meter() domain.IMeter {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IMeter)
}

// WithDatatable sets the datatable
func (a *AppMock) WithDatatable(datatable domain.IDatatable) domain.IApp {
	args := a.Called(datatable)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Datatable gets the datatable service
func (a *AppMock) Datatable() domain.IDatatable {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IDatatable)
}

// WithStateMachine sets the state machine service
func (a *AppMock) WithStateMachine(stateMachine instanceStateMachine.IStateMachineService) domain.IApp {
	args := a.Called(stateMachine)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// StateMachine gets the state machine service
func (a *AppMock) StateMachine() instanceStateMachine.IStateMachineService {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(instanceStateMachine.IStateMachineService)
}

// Config gets the app configurations
func (a *AppMock) Config() *appConfig.Config {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*appConfig.Config)
}

// ConfigFile gets the configuration file name
func (a *AppMock) ConfigFile() string {
	args := a.Called()
	return args.Get(0).(string)
}

// WithAdditionalConfigType sets an additional config type
func (a *AppMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := a.Called(obj)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (a *AppMock) Started() bool {
	args := a.Called()
	return args.Get(0).(bool)
}
