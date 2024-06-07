package main

import (
	"errors"
	"fmt"
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

func (p *Publisher2) SimpleSend(body []byte) error {

	// do something .........

	err := p.CustomSend(&amqp.Publishing{
		Headers:         nil,
		ContentType:     "application/octet-stream",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            body,
	})

	// do something .........

	return err
}

func (p *Publisher2) CustomSend(msg *amqp.Publishing) error {

	// do something .........

	ch := p.GetCh()
	if ch.IsClosed() {
		return errors.New("publisher send failed, because channel is closed")
	}
	err := ch.Publish(
		p.PublisherConfig.ExchangeName, // exchange
		p.PublisherConfig.RoutingKey,   // routing key
		p.PublisherConfig.Mandatory,    // mandatory
		p.PublisherConfig.Immediate,    // immediate
		*msg)
	if err != nil {
		return errors.New(fmt.Sprintf("publisher send failed, error: %v", err))
	}

	// do something .........

	return nil
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
		err = publisher.CustomSend(&amqp.Publishing{
			Headers:         nil,
			ContentType:     "",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            nil,
		})
		if err != nil {
			log.Println(err)
		}
	}
}
