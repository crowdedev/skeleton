package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
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

	go todos.NewServer(q.dispatcher).RegisterQueueConsumer()
}
