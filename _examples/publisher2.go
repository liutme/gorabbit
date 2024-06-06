package main

import (
	"github.com/liutme/gorabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const (
	ExchangeNamePublisher2 = "test-exchange-4"
	RoutingKeyPublisher2   = "queue-4"
)

type Publisher2 struct {
	gorabbit.Publisher
}

func (p *Publisher2) BuildPublisher() *gorabbit.PublisherConfig {
	p.PublisherConfig = gorabbit.PublisherConfig{
		ExchangeName: ExchangeNamePublisher2,
		RoutingKey:   RoutingKeyPublisher2,
	}
	return &p.PublisherConfig
}

func main() {
	publisher := &Publisher2{}

	rabbitClient := &gorabbit.Client{
		Config: gorabbit.ConnectionConfig{
			Host:     "127.0.0.1",
			Port:     "5672",
			UserName: "admin",
			Password: "admin",
			VHost:    "/",
		},
		Publishers: []gorabbit.IPublisher{
			publisher,
		},
	}

	rabbitClient.Init()
	for {
		time.Sleep(5 * time.Second)
		err := publisher.SimpleSend([]byte("a test message"))
		if err != nil {
			log.Println(err)
		}
		err = publisher.CustomSend(&amqp.Publishing{ // 可以使用Publishing内置参数自定义发送消息，没有限制
			ContentType: "text/plain",
			Body:        []byte("a test message"),
		})
		if err != nil {
			log.Println(err)
		}
	}
}
