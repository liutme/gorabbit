[English document](README.md)

------ 

# Go RabbitMQ 客户端库

这个库是对[Rabbitmq官方库amqp091-go](https://github.com/rabbitmq/amqp091-go)的封装,使用者无需关心连接和通道,只需要专注于自己的具体业务.

## 解决的问题

+ gorabbit实现了连接关闭重连。RabbitMQ团队维护的amqp091-go库并不关心连接状态，当按照简单示例注册客户端，可能不是一个安全的方案。

    + 维护一个RabbitMQ连接，并监听连接情况并进行全局缓存，当发生连接断开，会重新连接并更新到全局缓存中
    + 不同的操作角色会开启单独的通道(channel)，开启通道的同时也会监听它的状态。一旦关闭会重新打开通道。并通知到通道调用者，例如消费者（consumer）
  
+ gorabbit将生产者（publisher）进行了封装，使用者只需声明一些配置，则可以进行注册，目前支持发布简单消息和自定义消息

+ gorabbit将消费者（consumer）进行了封装，使用者只需声明一些配置，和实现一个监听函数，则可以进行注册。

  + 消费者监听接口函数，当接收到消息，则会回调到实现的监听函数中，只需要关注业务逻辑。

## 示例

可以查看[_examples](_examples)目录，里面有消费者和生产者示例，也可以查看下面的说明。也可以提交issue。

## 用法

### 客户端说明
gorabbit.Client是统一集成的客户端配置结构，包含了客户端连接配置，角色声明配置，具体结构如下。 
```go
type Client struct {
    Config     ConnectionConfig // 连接配置
    Consumers  []IConsumer      // 消费者列表
    Publishers []IPublisher     // 生产者列表
}
```

### 连接

使用时需要声明连接配置提供给gorabbit，用于连接和重连。

```go
rabbitClient := &gorabbit.Client{
    Config: gorabbit.ConnectionConfig{
        Host:     "127.0.0.1",
        Port:     "5672",
        UserName: "admin",
        Password: "admin",
        VHost:    "/",
    },
}
```


------

### 消费者:

1. **gorabbit.Consumer说明**

```go
type Consumer struct {
	Queue          Queue          // 队列配置
	ConsumerConfig ConsumerConfig // 消费者配置
	QueueBinding   QueueBinding   // 队列绑定配置
}
```

完整的配置集
```go
gorabbit.Consumer{
    Queue: gorabbit.Queue{
        Name:       "",    // 队列名称
        Durable:    false, // 是否持久化
        AutoDelete: false, // 是否自动删除
        Exclusive:  false, // 是否排他性
        NoWait:     false, // 是否等待RabbitMQ server确认
        Args:       nil,   // 其他自定义参数
    },
    ConsumerConfig: gorabbit.ConsumerConfig{
        Tag:       "",    // 消费者标签
        AutoAck:   false, // 是否自动确认
        Exclusive: false, // 是否排他性
        NoLocal:   false, // 可以不管
        NoWait:    false, // 是否等待RabbitMQ server确认
        Args:      nil,   // 其他自定义参数
    },
    QueueBinding: gorabbit.QueueBinding{
        Exchange: gorabbit.Exchange{ // 交换机配置，自动创建
            Name:       "",    // 交换机名称
            Kind:       "",    // 交换机类型
            Durable:    false, // 是否持久化
            AutoDelete: false, // 是否自动删除
            Internal:   false, // 是否为内部交换机，不用管
            NoWait:     false, // 是否等待RabbitMQ server确认
            Args:       nil,   // 其他自定义参数
        },
        RoutingKey: []string{QueueName}, // routingkey
	},
}
```

2. **注册一个简单的消费者示例**

   1. 声明一个消费者结构，并继承gorabbit.Consumer
   
   2. 实现Listener函数
   
      ```go
      type ExampleConsumer struct {
           gorabbit.Consumer
      }
    
       func (c *ExampleConsumer) Listener(delivery *amqp.Delivery) {
           body := string(delivery.Body)
           log.Printf("a message was received：%s", body)
       }
      ```
   3. 然后注册到客户端中

       ```go
       rabbitClient := &gorabbit.Client{
               Config: gorabbit.ConnectionConfig{
                   Host:     "127.0.0.1",
                   Port:     "5672",
                   UserName: "admin",
                   Password: "admin",
                   VHost:    "/",
               },
               Consumers: []gorabbit.IConsumer{
                   &ExampleConsumer{
                       Consumer: gorabbit.Consumer{
                           Queue: gorabbit.Queue{
                               Name:       "example-queue-1",
                               Durable:    true,
                               AutoDelete: false,
                           },
                           ConsumerConfig: gorabbit.ConsumerConfig{
                               AutoAck: true,
                           },
                           QueueBinding: gorabbit.QueueBinding{
                               Exchange: gorabbit.Exchange{
                                   Name:    "example-exchange-1",
                                   Kind:    "direct",
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
3. **另一种创建消费者的方式（可选）**  

   1. 声明一个消费者结构，并继承gorabbit.Consumer
   
   2. 实现Listener函数和BuildConsumer函数
   
       ```go
      type ExampleConsumer struct {
           gorabbit.Consumer
      }   
       
      func (c *ExampleConsumer) BuildConsumer() gorabbit.Consumer {
              c.Consumer = gorabbit.Consumer{
                  Queue: gorabbit.Queue{
                      Name:       QueueName2,
                      Durable:    true,
                      AutoDelete: false,
                  },
                  ConsumerConfig: gorabbit.ConsumerConfig{
                      AutoAck: true,
                  },
                  QueueBinding: gorabbit.QueueBinding{
                      Exchange: gorabbit.Exchange{
                          Name:    ExchangeName2,
                          Kind:    "direct",
                          Durable: true,
                      },
                      RoutingKey: []string{QueueName2},
                      NoWait:     false,
                  },
              }
              return c.Consumer
          }
      
          func (c *ExampleConsumer) Listener(delivery *amqp.Delivery) {
              body := string(delivery.Body)
              log.Printf("a message was received：%s, ConsumeConfig: %s", body, c.Queue.Name)
          }
       ```
   3. 然后注册到客户端中
      ```go
         rabbitClient := &gorabbit.Client{
             Config: gorabbit.ConnectionConfig{
                Host:     "127.0.0.1",
                Port:     "5672",
                UserName: "admin",
                Password: "admin",
                VHost:    "/",
            },
            Consumers: []gorabbit.IConsumer{
                &ExampleConsumer{},
            },
         }
    
          rabbitClient.Init()
      ``` 
------

### 生产者

1. **gorabbit.Publisher说明**

```go
    type Publisher struct {
        mtx             sync.RWMutex    // 读写锁
        ch              *amqp.Channel   // 连接通道
        PublisherConfig PublisherConfig // 生产者配置
    }
```

完整的配置集
```go
    gorabbit.PublisherConfig{
        ExchangeName: "", // 交换机名称
        RoutingKey:   "", // routingkey
        Mandatory:    false, // 对应的交换机和routingkey无法找到对应的queue时是否响应生产者
        Immediate:    false, // 消息发送后是否检测消费者存在，设置为true如果队列无消费者不会进入队列并返还消息
    }
```

2. **一个简单的生产者示例**

   1. 声明一个生产者结构，并继承gorabbit.Publisher
   
   ```go
     type ExamplePublisher struct {
	        gorabbit.Publisher
     }
   ```
    
   2. 然后注册到客户端中
   
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


3. **另一种创建生产者的方式（可选）**
   1. 声明一个生产者结构，并继承gorabbit.Publisher
   2. 实现BuildPublisher函数
   ```go
    type ExamplePublisher struct {
        gorabbit.Publisher
    }
    
    func (p *ExamplePublisher) BuildPublisher() *gorabbit.PublisherConfig {
        p.PublisherConfig = PublisherConfig: gorabbit.PublisherConfig{
                ExchangeName: "",
                RoutingKey:   "",
                Mandatory:    false,
                Immediate:    false,
            },
        return &p.PublisherConfig
    }
   ```
   3. 然后注册到客户端中

   ```go
    publisher := &ExamplePublisher{}
   
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

4. **生产者的消息发送**

   目前生产者发送消息的函数只有两个
   
   1. publisher.SimpleSend，简单的消息投递，参数都是默认amqp默认，只需要传入Body
   2. publisher.CustomSend，自定义消息投递，使用amqp.Publishing作为参数，可以自定义投递参数

    ```go
    publisher.SimpleSend([]byte("a test message"))
    
    publisher.CustomSend(&amqp.Publishing{
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
    ```

5. **重写发送函数**

发送函数支持重写，支持使用者自定义

```go
func (p *ExamplePublisher) SimpleSend(body []byte) error {

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

func (p *ExamplePublisher) CustomSend(msg *amqp.Publishing) error {
   
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

```
