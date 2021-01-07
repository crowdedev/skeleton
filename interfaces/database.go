package interfaces

import (
	configs "github.com/crowdeco/todo-service/configs"
	todos "github.com/crowdeco/todo-service/todos"
)

type (
	Database struct{}
)

func NewDatabase() configs.Application {
	return &Database{}
}

func (d *Database) Run() {
	todos.RegisterAutoMigration()
}
