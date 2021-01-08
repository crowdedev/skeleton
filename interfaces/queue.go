package interfaces

import (
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
	todos.NewTodo().Consume()
}
