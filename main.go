package main

import (
	configs "github.com/crowdeco/skeleton/configs"
	"github.com/crowdeco/skeleton/events"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	"github.com/crowdeco/skeleton/todos/listeners"
)

func init() {
	configs.LoadConfigs()
	configs.Env.ServiceName = "skeleton"
	configs.Env.Version = "v1.1@dev"
}

func main() {
	dispatcher := events.NewDispatcher()
	dispatcher.Register([]events.Listener{
		listeners.NewTodoSearch(),
	})

	database := interfaces.NewDatabase(dispatcher)
	go database.Run()

	grpc := interfaces.NewGRpc(dispatcher)
	go grpc.Run()

	queue := interfaces.NewQueue(dispatcher)
	go queue.Run()

	rest := interfaces.NewRest()
	rest.Run()
}
