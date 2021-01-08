package interfaces

import (
	configs "github.com/crowdeco/skeleton/configs"
	todos "github.com/crowdeco/skeleton/todos"
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
