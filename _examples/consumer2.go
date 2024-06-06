package main

import (
	"github.com/liutme/gorabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consume2 struct {
	gorabbit.Consumer
}

const (
	QueueName2    = "queue-2"
	ExchangeName2 = "exchange-2"
)

func (c *Consume2) BuildConsumer() gorabbit.Consumer {
	c.Consumer = gorabbit.Consumer{
		Queue: gorabbit.Queue{
			Name:       QueueName2,
			Durable:    true,
			AutoDelete: false,
		},
		ConsumerConfig: gorabbit.ConsumerConfig{
			AutoAck: true,
		},
		QueueBinding: gorabbit.QueueBinding{
			Exchange: gorabbit.Exchange{
				Name:    ExchangeName2,
				Kind:    "direct",
				Durable: true,
			},
			RoutingKey: []string{QueueName2},
			NoWait:     false,
		},
	}
	return c.Consumer
}

func (c *Consume2) Listener(delivery *amqp.Delivery) {
	body := string(delivery.Body)
	log.Printf("a message was received：%s, ConsumeConfig: %s", body, c.Queue.Name)
}

func main() {
	rabbitClient := &gorabbit.Client{
		Config: gorabbit.ConnectionConfig{
			Host:     "127.0.0.1",
			Port:     "5672",
			UserName: "admin",
			Password: "admin",
			VHost:    "/",
		},
		Consumers: []gorabbit.IConsumer{
			&Consume2{},
		},
	}

	rabbitClient.Init()
	forever := make(chan bool)
	<-forever
}
