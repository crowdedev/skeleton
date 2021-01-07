package interfaces

import (
	configs "github.com/crowdeco/todo-service/configs"
	todos "github.com/crowdeco/todo-service/todos"
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
