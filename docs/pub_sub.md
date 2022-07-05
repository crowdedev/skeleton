# Implement Pub/Sub

- Add RabbitMQ config to `dics/container.go`

```go
{
    Name:  "bima:messenger:config",
    Scope: bima.Application,
    Build: func(dsn string) (amqp.Config, error) {
        return amqp.NewDurableQueueConfig(dsn), nil
    },
    Params: dingo.Params{
        "0": "amqp://guest:guest@localhost:5672",
    },
},
{
    Name:  "bima:publisher",
    Scope: bima.Application,
    Build: func(env *configs.Env, config amqp.Config) (*amqp.Publisher, error) {
        publisher, err := amqp.NewPublisher(config, watermill.NewStdLogger(env.Debug, env.Debug))
        if err != nil {
            return nil, nil
        }

        return publisher, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:messenger:config"),
    },
},
{
    Name:  "bima:consumer",
    Scope: bima.Application,
    Build: func(env *configs.Env, config amqp.Config) (*amqp.Subscriber, error) {
        consumer, err := amqp.NewSubscriber(config, watermill.NewStdLogger(env.Debug, env.Debug))
        if err != nil {
            return nil, nil
        }

        return consumer, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:messenger:config"),
    },
},
{
    Name:  "bima:messenger",
    Scope: bima.Application,
    Build: func(
        env *configs.Env,
        publisher *amqp.Publisher,
        consumer *amqp.Subscriber,
    ) (*messengers.Messenger, error) {
        if consumer == nil || publisher == nil {
            return nil, nil
        }

        color.New(color.FgCyan, color.Bold).Print("✓ ")
        fmt.Println("Pub/Sub configured")

        return messengers.New(env.Debug, publisher, consumer), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:publisher"),
        "2": dingo.Service("bima:consumer"),
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
