package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	todos "github.com/crowdeco/skeleton/todos"
)

type (
	Queue struct{}
)

func NewQueue() configs.Application {
	return &Queue{}
}

func (q *Queue) Run() {
	log.Printf("Starting Queue Consumer")

	todos.NewServer().RegisterQueueConsumer()
}
