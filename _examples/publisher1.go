package main

import (
	"github.com/liutme/gorabbit"
	"log"
	"time"
)

type Publisher1 struct {
	gorabbit.Publisher
}

const (
	ExchangeNamePublisher1 = "test-exchange-3"
	RoutingKeyPublisher1   = "queue-3"
)

func main() {
	publisher1 := &Publisher1{
		Publisher: gorabbit.Publisher{
			PublisherConfig: gorabbit.PublisherConfig{
				ExchangeName: ExchangeNamePublisher1,
				RoutingKey:   RoutingKeyPublisher1,
			},
		},
	}

	rabbitClient := &gorabbit.Client{
		Config: gorabbit.ConnectionConfig{
			Host:     "127.0.0.1",
			Port:     "5672",
			UserName: "admin",
			Password: "admin",
			VHost:    "/",
		},
		Publishers: []gorabbit.IPublisher{
			publisher1,
		},
	}

	rabbitClient.Init()
	for {
		time.Sleep(5 * time.Second)
		err := publisher1.SimpleSend([]byte("a test message"))
		if err != nil {
			log.Println(err)
		}
	}
}
