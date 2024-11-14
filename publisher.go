package gorabbit

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
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

func (c *Client) RegisterPublisher() {
	if len(c.Publishers) < 1 {
		c.Logger.Warn("no publishers registered")
		return
	}
	for _, publisher := range c.Publishers {
		channels := openChannel(c)
		p := publisher.BuildPublisher()
		ch := <-channels
		if ch.IsClosed() {
			c.Logger.fatal("publisher registration failed, because the channel listening to the MQ connection channel has been closed")
		}
		publisher.setCh(ch)
		go func() {
			for {
				select {
				case ch, ok := <-channels:
					if !ok {
						c.Logger.fatal("publisher stopped, because the channel listening to the MQ connection channel has been closed")
					}
					if ch.IsClosed() {
						continue
					}
					publisher.setCh(ch)
					c.Logger.Info("The publisher registered successfully", "RoutingKey", p.RoutingKey)
				}
			}
		}()
	}
}
