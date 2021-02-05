package events

import (
	"fmt"

	"github.com/crowdeco/skeleton/configs"
)

type Dispatcher struct {
	Events map[string]configs.Listener
}

func (d *Dispatcher) Register(listeners []configs.Listener) {
	for _, listener := range listeners {
		if _, ok := d.Events[listener.Listen()]; ok {
			panic(fmt.Sprintf("the '%s' event is already registered", listener.Listen()))
		}

		d.Events[listener.Listen()] = listener
	}
}

func (d *Dispatcher) Dispatch(name string, event interface{}) error {
	if _, ok := d.Events[name]; !ok {
		return fmt.Errorf("the '%s' event is already registered", name)
	}

	d.Events[name].Handle(event)

	return nil
}
