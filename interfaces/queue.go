package interfaces

import (
	"log"

	events "github.com/crowdeco/skeleton/events"
	todos "github.com/crowdeco/skeleton/todos"
)

type Queue struct {
	Dispatcher *events.Dispatcher
}

func (q *Queue) Run() {
	log.Printf("Starting Queue Consumer")

	go todos.NewServer(q.Dispatcher).RegisterQueueConsumer()
}
