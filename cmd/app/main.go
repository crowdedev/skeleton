package main

import (
	"fmt"

	dic "github.com/crowdeco/skeleton/generated/dic"
)

func main() {
	container, _ := dic.NewContainer()

	database, err := container.SafeGetCoreInterfaceDatabase()
	if err != nil {
		fmt.Errorf("Error Database: %s", err.Error())
		return
	}

	go database.Run()

	grpc, err := container.SafeGetCoreInterfaceGrpc()
	if err != nil {
		fmt.Errorf("Error gRPC: %s", err.Error())
		return
	}

	go grpc.Run()

	queue, err := container.SafeGetCoreInterfaceQueue()
	if err != nil {
		fmt.Errorf("Error Queue: %s", err.Error())
		return
	}

	go queue.Run()

	rest, err := container.SafeGetCoreInterfaceRest()
	if err != nil {
		fmt.Errorf("Error REST: %s", err.Error())
		return
	}

	rest.Run()
}
