package gorabbit

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type IPublisher interface {
	BuildPublisher() *PublisherConfig
	setCh(ch *amqp.Channel)
	GetCh() *amqp.Channel
	SimpleSend([]byte) error
	CustomSend(msg *amqp.Publishing) error
}

type Publisher struct {
	mtx             sync.RWMutex
	ch              *amqp.Channel
	PublisherConfig PublisherConfig
}

func (p *Publisher) BuildPublisher() *PublisherConfig {
	return &p.PublisherConfig
}

func (p *Publisher) setCh(ch *amqp.Channel) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.ch = ch
}

func (p *Publisher) GetCh() *amqp.Channel {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	return p.ch
}

func (p *Publisher) SimpleSend(body []byte) error {
	return p.CustomSend(&amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        body,
	})
}

func (p *Publisher) CustomSend(msg *amqp.Publishing) error {
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
	return nil
}

func PublisherRegisters(publisher ...IPublisher) {
	for _, p := range publisher {
		publisherRegister(p)
	}
}

func publisherRegister(ps IPublisher) {
	channels := openChannel()
	p := ps.BuildPublisher()
	ch := <-channels
	if ch.IsClosed() {
		log.Fatal("publisher registration failed, because the channel listening to the MQ connection channel has been closed")
	}
	ps.setCh(ch)
	go func() {
		for {
			select {
			case ch, ok := <-channels:
				if !ok {
					log.Fatal("publisher stopped, because the channel listening to the MQ connection channel has been closed")
				}
				if ch.IsClosed() {
					continue
				}
				ps.setCh(ch)
				log.Printf("The publisher registered successfully, RoutingKey: %s", p.RoutingKey)
			}
		}
	}()
}
