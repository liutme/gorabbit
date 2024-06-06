package main

import (
	"github.com/liutme/gorabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consume1 struct {
	gorabbit.Consumer
}

func (c *Consume1) Listener(delivery *amqp.Delivery) {
	body := string(delivery.Body)
	log.Printf("a message was received：%s", body)
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
			&Consume1{
				Consumer: gorabbit.Consumer{
					Queue: gorabbit.Queue{
						Name:       "",
						Durable:    false,
						AutoDelete: false,
						Exclusive:  false,
						NoWait:     false,
						Args:       nil,
					},
					ConsumerConfig: gorabbit.ConsumerConfig{
						Tag:       "",
						AutoAck:   false,
						Exclusive: false,
						NoLocal:   false,
						NoWait:    false,
						Args:      nil,
					},
					QueueBinding: gorabbit.QueueBinding{
						Exchange: gorabbit.Exchange{
							Name:       "",
							Kind:       "",
							Durable:    false,
							AutoDelete: false,
							Internal:   false,
							NoWait:     false,
							Args:       nil,
						},
						RoutingKey: []string{""},
					},
				},
			},
		},
	}

	rabbitClient.Init()

	forever := make(chan bool)
	<-forever
}
