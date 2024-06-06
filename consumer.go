package gorabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

// IConsumer interface, including a configuration builder function and a consumption listener function.
type IConsumer interface {
	BuildConsumer() Consumer
	Listener(delivery *amqp.Delivery)
}
type Consumer struct {
	Queue          Queue
	ConsumerConfig ConsumerConfig
	QueueBinding   QueueBinding
}

// BuildConsumer A built-in consumer builder that directly returns the consumer configuration defined in the main thread.
func (c Consumer) BuildConsumer() Consumer {
	return c
}

/*
ConsumerRegisters This function registers multiple consumers
*/
func ConsumerRegisters(cs ...IConsumer) {
	for _, c := range cs {
		consumerRegister(c)
	}
}

/*
consumerRegister This function registers a single consumer
a connection channel will be opened, then the queues, exchanges you defined will be created,
and the queues will be bound to the exchanges.
After these operations are completed, a consumer listener will be registered,
which is able to receive messages and pass them to the listener implementation you declared.
*/
func consumerRegister(c IConsumer) {
	channels := openChannel()
	configs := c.BuildConsumer()
	go func() {
		var deliveries <-chan amqp.Delivery
		for {
			select {
			case ch, ok := <-channels:
				if !ok {
					log.Fatalf("ConsumerConfig stoped, because because the channel listening to the MQ connection channel has been closed, queue: %s", configs.Queue.Name)
				}
				if ch.IsClosed() {
					continue
				}
				deliveries = configs.createConsumer(ch)
				log.Printf("The consumer registered successfully, queue: %s", configs.Queue.Name)
			case d, ok := <-deliveries:
				if !ok {
					log.Printf("ConsumerConfig stoped, queue: %s, waiting for the rebuild rebuild...", configs.Queue.Name)
					time.Sleep(3 * time.Second)
					continue
				}
				c.Listener(&d)
			}
		}
	}()
}

// createConsumer Based on your configuration, create a consumer.
func (c Consumer) createConsumer(ch *amqp.Channel) <-chan amqp.Delivery {
	_, err := ch.QueueDeclare(
		c.Queue.Name,       // Name
		c.Queue.Durable,    // Durable
		c.Queue.AutoDelete, // delete when unused
		c.Queue.Exclusive,  // exclusive
		c.Queue.NoWait,     // no-wait
		c.Queue.Args,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to queue declare, queue: %s, error: %s", c.Queue.Name, err)
	}
	if len(c.QueueBinding.Exchange.Name) > 1 {
		err = ch.ExchangeDeclare(
			c.QueueBinding.Exchange.Name,       // Name
			c.QueueBinding.Exchange.Kind,       // type
			c.QueueBinding.Exchange.Durable,    // Durable
			c.QueueBinding.Exchange.AutoDelete, // auto-deleted
			c.QueueBinding.Exchange.Internal,   // Internal
			c.QueueBinding.Exchange.NoWait,     // no-wait
			c.QueueBinding.Exchange.Args,       // arguments
		)
		if err != nil {
			log.Fatalf("Failed to exchange declare, exchange: %s, error: %s", c.QueueBinding.Exchange.Name, err)
		}
		if len(c.QueueBinding.RoutingKey) < 1 {
			err = ch.QueueBind(
				c.Queue.Name,                 // queue Name
				"",                           // routing key
				c.QueueBinding.Exchange.Name, // exchange
				c.QueueBinding.NoWait,
				c.QueueBinding.Args)
		} else {
			for _, rk := range c.QueueBinding.RoutingKey {
				err = ch.QueueBind(
					c.Queue.Name,                 // queue Name
					rk,                           // routing key
					c.QueueBinding.Exchange.Name, // exchange
					c.QueueBinding.NoWait,
					c.QueueBinding.Args)
			}
		}
		if err != nil {
			log.Fatalf("Failed to queue binding: %s", err)
		}
	}
	deliveries, err := ch.Consume(
		c.Queue.Name,
		c.ConsumerConfig.Tag,
		c.ConsumerConfig.AutoAck,
		c.ConsumerConfig.Exclusive,
		c.ConsumerConfig.NoLocal,
		c.ConsumerConfig.NoWait,
		c.ConsumerConfig.Args,
	)
	return deliveries
}
