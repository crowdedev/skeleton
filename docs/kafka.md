# Use Kafka as Message Broker

- Add Kafka config to `dics/container.go`

```go
{
    Name:  "bima:kafka:publisher",
    Scope: bima.Application,
    Build: func(env *configs.Env, hosts []string) (*kafka.Publisher, error) {
        publisher, err := kafka.NewPublisher(kafka.PublisherConfig{
            Brokers:   hosts,
            Marshaler: kafka.DefaultMarshaler{},
        }, watermill.NewStdLogger(env.Debug, env.Debug))
        if err != nil {
            return nil, nil
        }

        return publisher, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": []string{"kafka:9092"},
    },
},
{
    Name:  "bima:kafka:consumer",
    Scope: bima.Application,
    Build: func(env *configs.Env, hosts []string, consumerGroup string) (*kafka.Subscriber, error) {
        saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	    saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
        consumer, err := kafka.NewSubscriber(kafka.SubscriberConfig{
            Brokers:               hosts,
            Unmarshaler:           kafka.DefaultMarshaler{},
            OverwriteSaramaConfig: saramaSubscriberConfig,
            ConsumerGroup:         consumerGroup,
        }, watermill.NewStdLogger(env.Debug, env.Debug))
        if err != nil {
            return nil, nil
        }

        return consumer, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": []string{"kafka:9092"},
        "2": "consumer_group"
    },
},
{
    Name:  "bima:kafka:broker",
    Scope: bima.Application,
    Build: func(publisher *kafka.Publisher, consumer *kafka.Subscriber) (messengers.Broker, error) {
        return brokers.NewKafka(publisher, consumer), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:kafka:publisher"),
        "1": dingo.Service("bima:kafka:consumer"),
    },
},
{
    Name:  "bima:messenger",
    Scope: bima.Application,
    Build: func(
        env *configs.Env,
        broker messengers.Broker,
    ) (*messengers.Messenger, error) {
        if consumer == nil || publisher == nil {
            return nil, nil
        }

        color.New(color.FgCyan, color.Bold).Print("✓ ")
        fmt.Println("Pub/Sub configured")

        return messengers.New(env.Debug, broker), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:kafka:broker"),
    },
},
```

- Register consumer server to `dics/container.go`, the name must `bima:interface:consumer`

```go
{
    Name:  "bima:interface:consumer",
    Scope: bima.Application,
    Build: (*interfaces.Queue)(nil),
    Params: dingo.Params{
        "Messenger": dingo.Service("bima:messenger"),
    },
},
```

## Consumer

To consume some message, you just need to override `Consume()` in your `server.go` like below

```go
func (s *Server) Consume(messenger *messengers.Messenger) {
	messages, err := messenger.Consume("topic")
	if err != nil {
		return
	}

	for _, message := range messages {
		//Do with message
	}
}
```

## Publisher

To publish some message in your module, need to add `Messenger` in your `module.go`

```go
type Module struct {
	*bima.Module
	Model     *Todo
	Messenger *messengers.Messenger
	grpcs.UnimplementedTodosServer
}
```

and then change `dic.go` and add `bima:messenger` as params

```go
{
    Name:  "module:todo",
    Scope: bima.Application,
    Build: (*Module)(nil),
    Params: dingo.Params{
        "Model":     dingo.Service("module:todo:model"),
        "Module":    dingo.Service("bima:module"),
        "Messenger": dingo.Service("bima:messenger"),
    },
},
```

now you can publish message using `m.Messenger.Publish()` function
