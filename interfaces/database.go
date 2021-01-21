package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	todos "github.com/crowdeco/skeleton/todos"
)

type database struct{}

func NewDatabase() configs.Application {
	return &database{}
}

func (d *database) Run() {
	log.Printf("Starting DB Auto Migration")

	todos.NewServer().RegisterAutoMigrate()
}
