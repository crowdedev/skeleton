package events

import (
	"fmt"
)

type Dispatcher struct {
	events map[string]Listener
}

func NewDispatcher(listeners []Listener) *Dispatcher {
	dispatcher := Dispatcher{
		events: make(map[string]Listener),
	}

	for _, listener := range listeners {
		if _, ok := dispatcher.events[listener.Listen()]; ok {
			panic(fmt.Sprintf("the '%s' event is already registered", listener.Listen()))
		}

		dispatcher.events[listener.Listen()] = listener
	}

	return &dispatcher
}

func (d *Dispatcher) Dispatch(name string, event interface{}) error {
	if _, ok := d.events[name]; !ok {
		return fmt.Errorf("the '%s' event is already registered", name)
	}

	d.events[name].Handle(event)

	return nil
}
