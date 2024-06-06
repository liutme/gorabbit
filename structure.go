package gorabbit

import amqp "github.com/rabbitmq/amqp091-go"

// ConnectionConfig connection configuration parameters
type ConnectionConfig struct {
	Host     string
	Port     string
	UserName string
	Password string
	VHost    string
}

// Queue queue configuration parameters
type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// QueueBinding queueBinding configuration parameters
type QueueBinding struct {
	Exchange   Exchange
	RoutingKey []string
	NoWait     bool
	Args       amqp.Table
}

// Exchange exchange configuration parameters
type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// ConsumerConfig consumer configuration parameters
type ConsumerConfig struct {
	Tag       string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// PublisherConfig publisher configuration parameters
type PublisherConfig struct {
	ExchangeName string
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
}
