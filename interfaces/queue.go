package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	parents "github.com/crowdeco/skeleton/parents"
	todos "github.com/crowdeco/skeleton/todos"
)

type queue struct {
	dispatcher *events.Dispatcher
}

func NewQueue(dispatcher *events.Dispatcher) configs.Application {
	return &queue{dispatcher: dispatcher}
}

func (q *queue) Run() {
	log.Printf("Starting Queue Consumer")

	go parents.NewServer(q.dispatcher).RegisterQueueConsumer()
	go todos.NewServer(q.dispatcher).RegisterQueueConsumer()
}
