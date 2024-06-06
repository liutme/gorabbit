package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorabbit"
	"log"
)

type Consume1 struct {
	gorabbit.Consumer
}

func (c Consume1) Listener(delivery *amqp.Delivery) {
	body := string(delivery.Body)
	log.Printf("a message was received：%s", body)
}

const (
	QueueName    = "queue-1"
	ExchangeName = "exchange-1"
)

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
			Consume1{
				Consumer: gorabbit.Consumer{
					Queue: gorabbit.Queue{
						Name:       QueueName,
						Durable:    true,
						AutoDelete: false,
					},
					ConsumerConfig: gorabbit.ConsumerConfig{
						AutoAck: true,
					},
					QueueBinding: gorabbit.QueueBinding{
						Exchange: gorabbit.Exchange{
							Name:    ExchangeName,
							Kind:    "direct",
							Durable: true,
						},
						RoutingKey: []string{QueueName},
					},
				},
			},
		},
	}

	rabbitClient.Init()

	forever := make(chan bool)
	<-forever
}
