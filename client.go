package gorabbit

// Client Aggregation Configuration
type Client struct {
	Config     ConnectionConfig
	Consumers  []IConsumer
	Publishers []IPublisher
}

/*
Init Initialization function.
creating connections, consumer registration, and publisher registration.
*/
func (c *Client) Init() {

	c.createConnections()

	if len(c.Consumers) > 0 {
		ConsumerRegisters(c.Consumers...)
	}
	if len(c.Publishers) > 0 {
		PublisherRegisters(c.Publishers...)
	}
}
