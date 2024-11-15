package rabbitmq

// / Configuration structs
type ExchangeConfig struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Durable bool   `json:"durable"`
}

type QueueConfig struct {
	Name    string `json:"name"`
	Durable bool   `json:"durable"`
}

type BindingConfig struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	RoutingKey  string `json:"routing_key"`
}

type RabbitMQConfig struct {
	Exchanges []ExchangeConfig `json:"exchanges"`
	Queues    []QueueConfig    `json:"queues"`
	Bindings  []BindingConfig  `json:"bindings"`
}
