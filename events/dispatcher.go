package events

import (
	"errors"
	"sort"

	"github.com/crowdeco/skeleton/configs"
)

type Dispatcher struct {
	Events map[string][]configs.Listener
}

func (d *Dispatcher) Register(listeners []configs.Listener) {
	d.Events = make(map[string][]configs.Listener)
	sort.Slice(listeners, func(i, j int) bool {
		return listeners[i].Priority() > listeners[j].Priority()
	})

	for _, listener := range listeners {
		if _, ok := d.Events[listener.Listen()]; !ok {
			d.Events[listener.Listen()] = []configs.Listener{}
		}

		d.Events[listener.Listen()] = append(d.Events[listener.Listen()], listener)
	}
}

func (d *Dispatcher) Dispatch(name string, event interface{}) error {
	if _, ok := d.Events[name]; !ok {
		return errors.New("Unregistered event")
	}

	for _, listener := range d.Events[name] {
		listener.Handle(event)
	}

	return nil
}
