package main

import (
	configs "github.com/crowdeco/skeleton/configs"
	dic "github.com/crowdeco/skeleton/dics/generated/dic"
)

func init() {
	configs.LoadConfigs()
	configs.Env.ServiceName = "skeleton"
	configs.Env.Version = "v1.1@dev"
}

func main() {
	container, _ := dic.NewContainer()

	database := container.GetCoreInterfaceDatabase()
	go database.Run()

	grpc := container.GetCoreInterfaceGrpc()
	go grpc.Run()

	queue := container.GetCoreInterfaceQueue()
	go queue.Run()

	rest := container.GetCoreInterfaceRest()
	rest.Run()
}
