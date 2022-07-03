# Implement Pub/Sub

## Config

Add `AMQP_HOST`, `AMQP_PORT`, `AMQP_USER` and `AMQP_PASSWORD` to `.env`

```bash
AMQP_HOST=localhost
AMQP_PORT=5672
AMQP_USER=guest
AMQP_PASSWORD=guest
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
    Build: (*Module)(nil),
    Params: dingo.Params{
        "Model":     dingo.Service("module:todo:model"),
        "Module":    dingo.Service("bima:module"),
        "Messenger": dingo.Service("bima:messenger"),
    },
}
```

now you can publish message using `m.Messenger.Publish()` function
