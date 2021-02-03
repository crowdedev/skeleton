package events

type Listener interface {
	Handle(event interface{})
	Listen() string
}
