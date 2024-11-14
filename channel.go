package gorabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

/*
openChannel This function initiates an MQ connection channel, which is then placed into a Golang channel and returned to the caller.
The caller can monitor this Golang channel to retrieve the MQ connection channel.
At the same time, the function manages the connection status of the MQ connection channel.
Once it detects that the MQ connection channel is closed, it will reopen a new MQ connection channel and resend it to the Golang channel,
allowing the caller to regain a healthy MQ connection channel without having to worry about its health status.
*/
func openChannel(c *Client) chan *amqp.Channel {
	conn := getConnection()
	if conn == nil || conn.IsClosed() {
		c.Logger.fatal("Failed to open a channel,because connection is nil or closed")
	}
	ch, err := conn.Channel()
	if err != nil {
		c.Logger.fatal("Failed to open a channel", "err", err)
	}

	channels := make(chan *amqp.Channel, 1)
	channels <- ch
	go func() {
		notifyClose := make(chan *amqp.Error)
		ch.NotifyClose(notifyClose)
		for {
			select {
			case closeErr := <-notifyClose:
				c.Logger.Warn("Channel closed,attempting to reopen channel...", "err", closeErr)
				time.Sleep(2 * time.Second) // 等待一段时间后重连
				conn = getConnection()
				if conn == nil || conn.IsClosed() {
					c.Logger.Warn("Failed to reopen a channel,because connection is nil or closed")
					continue
				}
				ch, err = conn.Channel()
				if err != nil {
					log.Printf(err.Error())
					continue
				}
				notifyClose = make(chan *amqp.Error)
				ch.NotifyClose(notifyClose)
				channels <- ch
				c.Logger.Warn("1 channel has been reopened")
			}
		}
	}()
	return channels
}
