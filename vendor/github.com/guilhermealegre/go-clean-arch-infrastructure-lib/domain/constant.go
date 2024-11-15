package domain

// Environments
const (
	// ProductionEnv production environment
	ProductionEnv = "prd"
	// LocalEnv local environment
	LocalEnv = "local"
)

// Constants
const (
	NotApplicable = "n/a"
)

// Sub Types
const (
	DefaultSubType  SubType = "service"
	DatabaseSubType SubType = "database"
	RabbitSubType   SubType = "rabbit"
	ElasticSubType  SubType = "elastic"
	RedisSubType    SubType = "redis"
	SQSSubType      SubType = "sqs"
)
