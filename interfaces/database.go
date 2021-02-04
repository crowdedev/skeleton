package interfaces

import (
	"log"

	events "github.com/crowdeco/skeleton/events"
	todos "github.com/crowdeco/skeleton/todos"
)

type Database struct {
	Dispatcher *events.Dispatcher
}

func (d *Database) Run() {
	log.Printf("Starting DB Auto Migration")

	todos.NewServer(d.Dispatcher).RegisterAutoMigrate()
}
