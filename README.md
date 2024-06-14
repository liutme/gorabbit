[中文文档](README_CN.md)

------ 

# Go RabbitMQ Client Library

[![Go Reference](https://pkg.go.dev/badge/github.com/liutme/gorabbit.svg)](https://pkg.go.dev/github.com/liutme/gorabbit)

This library is a wrapper of [Rabbitmq official library amqp091-go](https://github.com/rabbitmq/amqp091-go). Users do not need to care about connections and channels, but only need to focus on their specific business.

## Problems solved

+ gorabbit implements connection closing and reconnection. The amqp091-go library maintained by the RabbitMQ team does not care about the connection status. When registering a client according to a simple example, it may not be a safe solution.

    + Maintain a RabbitMQ connection, monitor the connection status and cache it globally. When a connection is disconnected, it will reconnect and update to the global cache.
    + Different operation roles will open separate channels, and while opening the channel, they will also monitor its status. Once closed, the channel will be reopened. And notify the channel caller, such as the consumer

+ gorabbit encapsulates the producer (publisher). Users only need to declare some configurations to register. Currently, it supports publishing simple messages and custom messages.

+ gorabbit encapsulates the consumer (consumer). Users only need to declare some configurations and implement a listening function to register.

    + The consumer listening interface function will be called back to the implemented listening function when a message is received. You only need to focus on the business logic.

## Example

You can view the [_examples](_examples) directory, which contains consumer and producer examples, or you can view the instructions below. You can also submit an issue.

## Usage

### Client Description gorabbit.Client
is a unified and integrated client configuration structure, which includes client connection configuration and role declaration configuration. The specific structure is as follows.
```go 
type Client struct { 
    Config ConnectionConfig // Connection configurationConsumers 
    []IConsumer // Consumer listPublishers 
    []IPublisher // Producer list 
} 
``` 

### Connection

When using, you need to declare the connection configuration to provide to gorabbit for connection and reconnection.

```go 
rabbitClient := &gorabbit.Client{ 
    Config: gorabbit.ConnectionConfig{ 
        Host: "127.0.0.1", 
        Port: "5672", 
        UserName: "admin", 
        Password: "admin", 
        VHost: "/", 
    }, 
} 
``` 


------ 

### Consumer:

1. **gorabbit.Consumer Description**

```go 
type Consumer struct { 
	Queue Queue // Queue configuration 
	ConsumerConfig ConsumerConfig // Consumer configuration 
	QueueBinding QueueBinding // Queue binding configuration 
} 
``` 

Complete configuration set
```go
gorabbit.Consumer{
    Queue: gorabbit.Queue{
        Name: "", // Queue name
        Durable: false, // Whether to persist
        AutoDelete: false, // Whether to automatically delete
        Exclusive: false, // Whether exclusive
        NoWait: false, // Whether to wait for RabbitMQ server confirmation
        Args: nil, // Other custom parameters
    },
    ConsumerConfig: gorabbit.ConsumerConfig{
        Tag: "", // Consumer tag
        AutoAck: false, // Automatic confirmation
        Exclusive: false, // Exclusive
        NoLocal: false, // No matter
        NoWait: false, // Wait for RabbitMQ server confirmation
        Args: nil, // Other custom parameters
    },
    QueueBinding: gorabbit.QueueBinding{
        Exchange: gorabbit.Exchange{ // Exchange configuration, automatic creation
            Name: "", // Exchange name
            Kind: "", // Exchange type
            Durable: false, // Whether to persist
            AutoDelete: false, // Whether to automatically delete
            Internal: false, // Is it an internal switch, don't worry about it
            NoWait: false, // Whether to wait for RabbitMQ server confirmation
            Args: nil, // Other custom parameters
        },
        RoutingKey: []string{QueueName}, // routingkey
    },
}
``` 

2. **Register a simple consumer example** 

   1. Declare a consumer structure and inherit gorabbit.Consumer 
   
   2. Implement the Listener function
   
        ```go
        type ExampleConsumer struct { 
           gorabbit.Consumer 
        } 
        
        func (c *ExampleConsumer) Listener(delivery *amqp.Delivery) { 
           body := string(delivery.Body) 
           log.Printf("a message was received：%s", body) 
        } 
        ``` 
   
   3. Then register to the client
        ```go
            rabbitClient := &gorabbit.Client{ 
               Config: gorabbit.ConnectionConfig{ 
                   Host: "127.0.0.1", 
                   Port: : "5672", 
                   UserName: "admin", 
                   Password: "admin", 
                   VHost: "/", 
               }, 
               Consumers: []gorabbit.IConsumer{ 
                   &ExampleConsumer{ 
                       Consumer: gorabbit.Consumer{ 
                           Queue: gorabbit.Queue{ 
                               Name: "example-queue-1", 
                               Durable: true, 
                               AutoDelete: false, 
                           }, 
                           ConsumerConfig: gorabbit.ConsumerConfig{ 
                               AutoAck: true, 
                           }, 
                           QueueBinding: gorabbit.QueueBinding{ 
                               Exchange: gorabbit.Exchange{ 
                                   Name: "example-exchange-1", 
                                   Kind: "direct", 
                                   Durable: true, 
                               }, 
                               RoutingKey: []string{"example-queue-1"}, 
                           }, 
                       }, 
                   }, 
               }, 
           } 
        
           rabbitClient.Init() 
        ``` 
3. **Another way to create a consumer (optional)**   

   1. Declare a consumer struct and inherit gorabbit.Consumer 
   
   2. Implement Listener function and BuildConsumer function 
   
       ```go 
      type ExampleConsumer struct { 
           gorabbit.Consumer
      }    
       
      func (c *ExampleConsumer) BuildConsumer() gorabbit.Consumer { 
              c.Consumer = gorabbit.Consumer{ 
                  Queue: gorabbit.Queue{ 
                      Name: QueueName2, 
                      Durable: true, 
                      AutoDelete: false, 
                  }, 
                  ConsumerConfig: gorabbit.ConsumerConfig{ 
                      AutoAck: true, 
                  }, 
                  QueueBinding: gorabbit.QueueBinding{ 
                      Exchange: gorabbit.Exchange{ 
                          Name: ExchangeName2, 
                          Kind: "direct", 
                          Durable: true, 
                      }, 
                      RoutingKey: []string{QueueName2}, 
                      NoWait: false, 
                  }, 
              } 
              return c.Consumer 
          } 
      
          func (c *ExampleConsumer) Listener(delivery *amqp.Delivery) { 
              body := string(delivery.Body) 
              log.Printf("a message was received: %s, ConsumeConfig: %s", body, c.Queue.Name) 
          } 
       ``` 
   3. Then register to the client
        ```go 
            rabbitClient 
             := &gorabbit.Client{ 
                 Config: gorabbit.ConnectionConfig{ 
                    Host: "127.0.0.1", 
                    Port: "5672", 
                    UserName: "admin", 
                    Password: "admin", 
                    VHost: "/", 
                }, 
                Consumers: []gorabbit.IConsumer{ 
                    &ExampleConsumer{}, 
                }, 
             } 
            
              rabbitClient.Init() 
        ``` 
------ 

### Producer 

1. **gorabbit.Publisher Description** 

```go 
    type Publisher struct { 
        mtx sync.RWMutex // read-write lock 
        ch *amqp.Channel // connection channel 
        PublisherConfig PublisherConfig // producer configuration 
    } 
``` 

Complete configuration set
```go
    gorabbit.PublisherConfig{
        ExchangeName: "", // exchange name 
		RoutingKey: "",   // routingkey
        Mandatory: false, // Whether to respond to the producer when the corresponding switch and routingkey cannot find the corresponding queue
        Immediate: false, // Whether to detect the existence of consumers after the message is sent, set to true if there is no consumer in the queue and the message will not be entered into the queue and returned
    }
``` 
2. **A simple producer example** 
   1. Declare a producer structure and inherit gorabbit.Publisher

        ```go
            type ExamplePublisher struct {
                gorabbit.Publisher 
            }
        ```
   
   2. Then register it in the client
       ```go
           publisher := &ExamplePublisher{
               Publisher: gorabbit.Publisher{
                   PublisherConfig: gorabbit.PublisherConfig{
                       ExchangeName: "",
                       RoutingKey:   "",
                       Mandatory:    false,
                       Immediate:    false,
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
           }
       ```



3. **Another way to create a producer (optional)**
    1. Declare a producer structure and inherit goribbit.Publisher
    2. Implement the BuildPublisher function
   
       ```go
            type ExamplePublisher struct {
                gorabbit.Publisher
            }
            
            func (p *ExamplePublisher) BuildPublisher() *gorabbit.PublisherConfig {
                p.PublisherConfig = PublisherConfig: gorabbit.PublisherConfig{
                        ExchangeName: "",
                        RoutingKey: "",
                        Mandatory: false,
                        Immediate: false,
                    },
                return &p.PublisherConfig
            }
       ```
       
    3. Then register to the client

       ```go
            publisher := &ExamplePublisher{}
           
            rabbitClient := &gorabbit.Client{
                Config: gorabbit.ConnectionConfig{
                    Host: "127.0.0.1",
                    Port: "5672",
                    UserName: "admin",
                    Password: "admin",
                    VHost: "/",
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
            }
       ```

4. **Producer message sending**

   Currently, there are only two functions for producers to send messages:

    1. publisher.SimpleSend, simple message delivery, all parameters are default amqp default, only need to pass in Body
    2. publisher.CustomSend, custom message delivery, use amqp.Publishing as a parameter, you can customize the delivery parameters

        ```go
        publisher.SimpleSend([]byte("a test message"))
        
        publisher.CustomSend(&amqp.Publishing{
            Headers: nil,
            ContentType: "",
            ContentEncoding: "",
            DeliveryMode: 0,
            Priority: 0,
            CorrelationId: "",
            ReplyTo: "",
            Expiration: "",
            MessageId: "",
            Timestamp: time.Time{},
            Type: "",
            UserId: "",
            AppId: "",
            Body: nil,
        })
        ```

5. **Rewrite the send function**

    The sending function supports rewriting and user customization

```go
func (p *ExamplePublisher) SimpleSend(body []byte) error {

   // do something .........
   
   err := p.CustomSend(&amqp.Publishing{
      Headers: nil,
      ContentType: "application/octet-stream",
      ContentEncoding: "",
      DeliveryMode: 0,
      Priority: 0,
      CorrelationId: "",
      ReplyTo: "",
      Expiration: "",
      MessageId: "",
      Timestamp: time.Time{},
      Type: "",
      UserId: "",
      AppId: "",
      Body: body,
   })
   
   // do something .........
   
   return err
}

func (p *ExamplePublisher) CustomSend(msg *amqp.Publishing) error {
   
   // do something .........
   
   ch := p.GetCh()
   if ch.IsClosed() {
      return errors.New("publisher send failed, because channel is closed")
   }
   err := ch.Publish(
      p.PublisherConfig.ExchangeName, // exchange
      p.PublisherConfig.RoutingKey, // routing key
      p.PublisherConfig.Mandatory, // mandatory
      p.PublisherConfig.Immediate, // immediate
      *msg)
   if err != nil {
      return errors.New(fmt.Sprintf("publisher send failed, error: %v", err))
   }
   
   // do something .........
   
   return nil
}

```