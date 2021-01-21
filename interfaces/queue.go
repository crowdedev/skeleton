package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	todos "github.com/crowdeco/skeleton/todos"
)

type queue struct{}

func NewQueue() configs.Application {
	return &queue{}
}

func (q *queue) Run() {
	log.Printf("Starting Queue Consumer")

	go todos.NewServer().RegisterQueueConsumer()
}
