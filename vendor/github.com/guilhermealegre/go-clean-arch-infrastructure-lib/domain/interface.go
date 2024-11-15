package domain

import (
	"context"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-playground/validator/v10"
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/metric"

	meterConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/meter/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	s3Config "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/s3/config"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/datatable/database"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/datatable/elastic_search"

	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	msg "github.com/guilhermealegre/go-clean-arch-core-lib/pagination"
	appConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/app/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/database/config"
	elasticSearchConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/elastic_search/config"
	grpcConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/grpc/config"
	httpConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/http/config"
	loggerConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/config"
	rabbitmqConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/rabbitmq/config"
	redisConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/redis/config"
	sqsConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/sqs/config"
	stateMachineDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/state_machine/instance"
	tracerConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

// IApp App interface
type IApp interface {
	IService

	// Config gets the configuration
	Config() *appConfig.Config
	// ConfigFile the configuration file
	ConfigFile() string

	// WithLogger sets the logger
	WithLogger(logger ILogger) IApp
	// Logger gets the logger
	Logger() ILogger
	// WithDatabase sets the database
	WithDatabase(database IDatabase) IApp
	// Database gets the database
	Database() IDatabase
	// WithRabbitmq sets the rabbitmq
	WithRabbitmq(rabbitmq IRabbitMQ) IApp
	// Rabbitmq gets the rabbitmq
	Rabbitmq() IRabbitMQ
	// WithSQS sets the sqs
	WithSQS(sqs ISQS) IApp
	// SQS gets the sqs
	SQS() ISQS
	// WithRedis sets the redis
	WithRedis(redis IRedis) IApp
	// Redis gets the redis
	Redis() IRedis
	// WithElasticSearch sets the elastic search
	WithElasticSearch(elasticSearch IElasticSearch) IApp
	// ElasticSearch gets the elastic search
	ElasticSearch() IElasticSearch
	// WithGrpc sets the grpc
	WithGrpc(grpc IGrpc) IApp
	// Grpc gets the grpc
	Grpc() IGrpc
	// WithHttp sets the http
	WithHttp(http IHttp) IApp
	// Http gets the http
	Http() IHttp
	// WithValidator sets custom validators
	WithValidator(validator IValidator) IApp
	// Validator gets the validator
	Validator() IValidator
	// WithAws sets the aws connection
	WithAws(aws IAws) IApp
	// Aws gets the http
	Aws() IAws
	// WithTracer sets the tracer
	WithTracer(tracer ITracer) IApp
	// Tracer gets the tracer
	Tracer() ITracer
	// WithMeter sets the meter
	WithMeter(meter IMeter) IApp
	// Meter gets the meter
	Meter() IMeter
	// WithDatatable sets the datatable
	WithDatatable(datatable IDatatable) IApp
	// Datatable gets the datatable
	Datatable() IDatatable
	// WithAws sets the aws connection
	WithS3(s3 IS3) IApp
	// S3 gets the s3 Connection
	S3() IS3
	// WithStateMachine sets the state machine
	WithStateMachine(stateMachine stateMachineDomain.IStateMachineService) IApp
	// StateMachine gets the state machine service
	StateMachine() stateMachineDomain.IStateMachineService
	// WithAdditionalConfigType sets an additional config type
	WithAdditionalConfigType(obj interface{}) IApp
}

// IService service interface
type IService interface {
	// Name name of the service
	Name() string
	// Start starts the service
	Start() error
	// Stop stops the service
	Stop() error
	// Started true if service started
	Started() bool
}

// IMiddleware the interface of the middlewares
type IMiddleware interface {
	RegisterMiddlewares()
	GetHandlers() []gin.HandlerFunc
}

// IController the interface of the controllers
type IController interface {
	App() IApp
	Register()
	Json(ctx contextDomain.IContext, data interface{}, err ...error)
	JsonWithPagination(ctx contextDomain.IContext, data interface{}, pagination *msg.Pagination, err ...error)
}

// IHttp the interface for the http service
type IHttp interface {
	IService

	// ConfigFile gets the configuration file name
	ConfigFile() string
	// Config gets the configurations
	Config() *httpConfig.Config

	// WithController adds a controller
	WithMiddleware(controller IMiddleware) IHttp
	// WithController adds a controller
	WithController(controller IController) IHttp
	// WithRouter sets the router
	WithRouter(router *gin.Engine) IHttp
	// Router gets the router
	Router() *gin.Engine
}

// ILogger logger service interface
type ILogger interface {
	IService

	// ConfigFile gets the configuration file
	ConfigFile() string
	// Config gets the configurations
	Config() *loggerConfig.Config

	//Log the error
	Log() ILogging

	// Database Log
	DBLog(error) error
	// SQS Log
	SQSLog(error) error
	//Elastic Log
	ElasticLog(error) error
	//Redis Log
	RedisLog(error) error
}

// IDatabase database service interface
type IDatabase interface {
	IService

	// ConfigFile gets the configuration file
	ConfigFile() string
	// Config gets the configurations
	Config() *config.Config

	// Read the read connection
	Read() session.ISession
	// Write the write connection
	Write() session.ISession
}

// IElasticSearch elastic search service interface
type IElasticSearch interface {
	IService

	ConfigFile() string
	Config() *elasticSearchConfig.Config

	Client() *elasticsearch.Client
}

// IRabbitMQ rabbitmq service interface
type IRabbitMQ interface {
	IService

	// ConfigFile gets the configuration file
	ConfigFile() string
	// Config gets the configurations
	Config() *rabbitmqConfig.Config

	// Produce produces to the rabbitmq
	Produce(message any, exchange string, routingKey string) error
	// Consume consumes from the rabbitmq
	Consume(app IApp, queues string, handlers map[string]func(msg amqp.Delivery) bool)
	// WithConsumer adds a consumer to the rabbitmq
	WithConsumer(consumer IRabbitMQConsumer) IRabbitMQ
}

// ISQS sqs service interface
type ISQS interface {
	IService

	// WithAdditionalConfigType sets an additional config type
	WithAdditionalConfigType(obj interface{}) ISQS
	// ConfigFile gets the configuration file
	ConfigFile() string
	// Config gets the configurations
	Config() *sqsConfig.Config
	// WithConsumer adds a consumer to the rabbitmq
	WithConsumer(consumer ISQSConsumer) ISQS
	// Connection
	Connection(name string) ISQSConnection
}

type ISQSConnection interface {
	// Connect connect
	Connect() error
	// Produce produces to the sqs
	Produce(ctx context.Context, queue string, messageAttributes map[string]*sqs.MessageAttributeValue, messages ...string) error
	// Consume consumes from the sqs
	Consume(maskedQueue string, consumer ISQSConsumer)
}

// IS3 s3 service interface
type IS3 interface {
	IService

	// ConfigFile gets the configuration file
	ConfigFile() string
	// Config gets the configurations
	Config() *s3Config.Config
	// Client S3 client
	Client() *s3.Client
}

// IRedis redis service interface
type IRedis interface {
	IService

	// WithAdditionalConfigType sets an additional config type
	WithAdditionalConfigType(obj interface{}) IRedis

	// ConfigFile gets the configuration file
	ConfigFile() string
	// Config gets the configurations
	Config() *redisConfig.Config

	// Client gets the redis client
	Client() *redis.Client
}

// IGrpc grpc service interface
type IGrpc interface {
	IService

	InitServer() error
	InitClients() error

	// Config gets the configurations
	Config() *grpcConfig.Configs
	// ConfigFile gets the configuration file
	ConfigFile() string
	// GetClient gets the client by name
	GetClient(name string) (conn *grpc.ClientConn)
	// GetServer gets the server
	GetServer() (conn *grpc.Server, err error)
	// WithController adds a controller
	WithController(controller IController) IGrpc
}

// IConsumer defines the rabbitmq consumers interface
type IRabbitMQConsumer interface {
	// GetHandlers gets the handlers
	GetHandlers() map[string]func(msg amqp.Delivery) bool
	// GetQueue gets the queue
	GetQueue() string
}

type ISQSConsumer interface {
	// GetHandlers gets the handlers
	GetHandlers() map[string]func(ctx context.Context, msg *sqs.Message) bool
	// GetConnection gets the connection name
	GetConnection() string
	// GetQueue gets the queue
	GetQueue() string
	// GetAttributeNames gets the queue attribute names
	GetAttributeNames() []*string
	// GetMessageAttributeNames gets the message attribute names
	GetMessageAttributeNames() []*string
}

// ILogging defines what our logging lib need
type ILogging interface {
	// Log the error
	Do(err error, info ...*LoggerInfo)
	// Log Multi
	Multi(err []error, info ...*LoggerInfo)
	// Frontend
	Frontend(error string, level errors.Level, fe *Frontend)
	// Initialize
	Init(cgf loggerConfig.Config) ILogging
}

// IValidator Interface
type IValidator interface {
	IService
	// AddFieldValidators adds a custom field validator
	AddFieldValidators(v ...IFieldValidator) IValidator
	// AddStructValidators adds a custom struct validator
	AddStructValidators(v ...IStructValidator) IValidator
	// Validate validates the struct
	Validate(ctx contextDomain.IContext, v any) error
}

type IFieldValidator interface {
	Tag() string
	Func(a IApp) validator.FuncCtx
}

type IStructValidator interface {
	Struct() any
	Func(a IApp) validator.StructLevelFuncCtx
}

// IAWS interface
type IAws interface {
	IService
	// AddValidator adds a custom validator
	Connection() *aws.Config
}

// ITracer interface
type ITracer interface {
	IService

	// Config gets the configurations
	Config() *tracerConfig.Config
	// ConfigFile gets the configuration file
	ConfigFile() string
	// Trace traces data
	Trace(ctx context.Context, spanName string, data map[string]any, err error)
	// TraceCurrentSpan traces data to a current span
	TraceCurrentSpan(ctx context.Context, data map[string]any, err error)
}

// IMeter interface
type IMeter interface {
	IService

	// Config gets the configurations
	Config() *meterConfig.Config
	// ConfigFile gets the configuration file
	ConfigFile() string
	// Prometheus gets Prometheus meter
	Prometheus() metric.Meter
}

// IStateMachine interface
type IStateMachine interface {
	GetName() string
	GetStateMachine() stateMachineDomain.IStateMachine
	SetHandlers()
	AddStateMachineTrigger(stateMachine IStateMachine)
}

// IDatatable datatable service interface
type IDatatable interface {
	IService
	Database() database.IDatabase
	Elastic() elastic_search.IElastic
}

type FallbackReader interface {
	ReadLines() ([]string, error)
}

type FallbackWriter interface {
	io.Writer
	Remove() error
}
