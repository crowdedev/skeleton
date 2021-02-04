package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	parents "github.com/crowdeco/skeleton/parents"
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

	parents.NewServer(d.dispatcher).RegisterAutoMigrate()
	todos.NewServer(d.dispatcher).RegisterAutoMigrate()
}
