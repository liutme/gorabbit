package gorabbit

// Client Aggregation Configuration
type Client struct {
	Config     ConnectionConfig
	Consumers  []IConsumer
	Publishers []IPublisher
	Logger     ILogger
}

/*
Init Initialization function.
creating connections, consumer registration, and publisher registration.
*/
func (c *Client) Init() {
	if c.Logger == nil {
		c.Logger = Logger{}
	}

	c.createConnections()

	c.RegisterConsumer()

	c.RegisterPublisher()
}
