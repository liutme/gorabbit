package gorabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
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
func (consumer Consumer) BuildConsumer() Consumer {
	return consumer
}

/*
RegisterConsumer The function registers multiple consumers,
each of which will open a connection channel and then create the queue and exchange you defined, binding the queue to the exchange.
After completing these operations, the current consumer listener will be registered.
It can receive messages and pass them to the listener implementation you specified.
*/
func (c *Client) RegisterConsumer() {
	if len(c.Consumers) < 1 {
		c.Logger.Warn("no consumers registered")
		return
	}
	for _, consumer := range c.Consumers {
		channels := openChannel(c)
		configs := consumer.BuildConsumer()
		go func() {
			var deliveries <-chan amqp.Delivery
			for {
				select {
				case ch, ok := <-channels:
					if !ok {
						c.Logger.fatal("ConsumerConfig stoped, because because the channel listening to the MQ connection channel has been closed, queue: %s", configs.Queue.Name)
					}
					if ch.IsClosed() {
						continue
					}
					deliveries = configs.createConsumer(c, ch)
					c.Logger.Info("The consumer registered successfully, queue: %s", configs.Queue.Name)
				case d, ok := <-deliveries:
					if !ok {
						c.Logger.Warn("ConsumerConfig stoped, queue: %s, waiting for the rebuild rebuild...", configs.Queue.Name)
						time.Sleep(3 * time.Second)
						continue
					}
					consumer.Listener(&d)
				}
			}
		}()
	}
}

// createConsumer Based on your configuration, create a consumer.
func (consumer Consumer) createConsumer(c *Client, ch *amqp.Channel) <-chan amqp.Delivery {
	_, err := ch.QueueDeclare(
		consumer.Queue.Name,       // Name
		consumer.Queue.Durable,    // Durable
		consumer.Queue.AutoDelete, // delete when unused
		consumer.Queue.Exclusive,  // exclusive
		consumer.Queue.NoWait,     // no-wait
		consumer.Queue.Args,       // arguments
	)
	if err != nil {
		c.Logger.fatal("Failed to queue declare, queue: %s, error: %s", consumer.Queue.Name, err)
	}
	if len(consumer.QueueBinding.Exchange.Name) > 1 {
		err = ch.ExchangeDeclare(
			consumer.QueueBinding.Exchange.Name,       // Name
			consumer.QueueBinding.Exchange.Kind,       // type
			consumer.QueueBinding.Exchange.Durable,    // Durable
			consumer.QueueBinding.Exchange.AutoDelete, // auto-deleted
			consumer.QueueBinding.Exchange.Internal,   // Internal
			consumer.QueueBinding.Exchange.NoWait,     // no-wait
			consumer.QueueBinding.Exchange.Args,       // arguments
		)
		if err != nil {
			c.Logger.fatal("Failed to exchange declare, exchange: %s, error: %s", consumer.QueueBinding.Exchange.Name, err)
		}
		if len(consumer.QueueBinding.RoutingKey) < 1 {
			err = ch.QueueBind(
				consumer.Queue.Name,                 // queue Name
				"",                                  // routing key
				consumer.QueueBinding.Exchange.Name, // exchange
				consumer.QueueBinding.NoWait,
				consumer.QueueBinding.Args)
		} else {
			for _, rk := range consumer.QueueBinding.RoutingKey {
				err = ch.QueueBind(
					consumer.Queue.Name,                 // queue Name
					rk,                                  // routing key
					consumer.QueueBinding.Exchange.Name, // exchange
					consumer.QueueBinding.NoWait,
					consumer.QueueBinding.Args)
			}
		}
		if err != nil {
			c.Logger.fatal("Failed to queue binding: %s", err)
		}
	}
	deliveries, err := ch.Consume(
		consumer.Queue.Name,
		consumer.ConsumerConfig.Tag,
		consumer.ConsumerConfig.AutoAck,
		consumer.ConsumerConfig.Exclusive,
		consumer.ConsumerConfig.NoLocal,
		consumer.ConsumerConfig.NoWait,
		consumer.ConsumerConfig.Args,
	)
	return deliveries
}
