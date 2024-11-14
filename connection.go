package gorabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type Connection struct {
	mtx  sync.RWMutex
	conn *amqp.Connection
}

var connectionInstance *Connection

/*
*
createConnections This function is used to create an MQ connection, which is then cached globally after its creation.
Callers can retrieve the connection through this cache.
Meanwhile, the function maintains the status of the MQ connection.
Once it detects that the connection is closed, it will establish a new MQ connection and update the cache accordingly.
The new connection will be ready for use by callers.
TODO Currently, only a global single connection is supported, but in the future, it can be considered to support multiple connections.
*/
func (c *Client) createConnections() {
	// connect
	conn, err := amqp.Dial(c.Config.uri())
	if err != nil {
		c.Logger.fatal("Failed to connect to RabbitMQ", "err", err)
	}
	c.Logger.Info("Connected to RabbitMQ", "host", c.Config.Host)

	// reconnect
	go func() {
		notifyClose := make(chan *amqp.Error)
		conn.NotifyClose(notifyClose)
		for {
			select {
			case closeErr := <-notifyClose:
				c.Logger.Warn("Connection closed, attempting to reconnect...", "err", closeErr)
				time.Sleep(time.Second) // 等待一段时间后重连
				var err error
				conn, err = amqp.Dial(c.Config.uri())
				if err != nil {
					c.Logger.Warn("Failed to connect to RabbitMQ", "err", err)
					continue
				}
				notifyClose = make(chan *amqp.Error)
				conn.NotifyClose(notifyClose)
				setConnection(conn)
				c.Logger.Info("Reconnected to RabbitMQ")
			}
		}
	}()
	connectionInstance = &Connection{
		conn: conn,
	}
}

func (c *ConnectionConfig) uri() string {
	vhost := c.VHost
	if vhost == "/" {
		vhost = ""
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", c.UserName, c.Password, c.Host, c.Port, vhost)
}

func setConnection(conn *amqp.Connection) {
	connectionInstance.mtx.Lock()
	defer connectionInstance.mtx.Unlock()
	connectionInstance.conn = conn
}

func getConnection() *amqp.Connection {
	connectionInstance.mtx.RLock()
	defer connectionInstance.mtx.RUnlock()
	return connectionInstance.conn
}
