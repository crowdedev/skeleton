package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	todos "github.com/crowdeco/skeleton/todos"
)

type database struct {
	dispatcher *events.Dispatcher
}

func NewDatabase(dispatcher *events.Dispatcher) configs.Application {
	return &database{dispatcher: dispatcher}
}

func (d *database) Run() {
	log.Printf("Starting DB Auto Migration")

	todos.NewServer(d.dispatcher).RegisterAutoMigrate()
}
