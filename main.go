package main

import (
	configs "github.com/crowdeco/todo-service/configs"
	interfaces "github.com/crowdeco/todo-service/interfaces"
)

func init() {
	configs.LoadConfigs()
	configs.Env.ServiceName = "todo-service"
	configs.Env.Version = "v2.0@dev"
}

func main() {
	database := interfaces.NewDatabase()
	go database.Run()

	grpc := interfaces.NewGRpc()
	go grpc.Run()

	queue := interfaces.NewQueue()
	go queue.Run()

	rest := interfaces.NewRest()
	rest.Run()
}
