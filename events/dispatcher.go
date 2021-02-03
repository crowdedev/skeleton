package events

import (
	"fmt"
)

type Dispatcher struct {
	events map[string]Listener
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		events: make(map[string]Listener),
	}
}

func (d *Dispatcher) Register(listeners []Listener) error {
	for _, listener := range listeners {
		if _, ok := d.events[listener.Listen()]; ok {
			return fmt.Errorf("the '%s' event is already registered", listener.Listen())
		}

		d.events[listener.Listen()] = listener
	}

	return nil
}

func (d *Dispatcher) Dispatch(name string, event interface{}) error {
	if _, ok := d.events[name]; !ok {
		return fmt.Errorf("the '%s' event is already registered", name)
	}

	d.events[name].Handle(event)

	return nil
}
