package main

import (
	"fmt"

	"github.com/crowdeco/skeleton/configs"
	"github.com/crowdeco/skeleton/generated/dic"
)

func main() {
	container, _ := dic.NewContainer()

	var servers []configs.Server
	for _, m := range container.GetCoreConfigParser().Parse() {
		servers = append(servers, container.Get(fmt.Sprintf("%s:server", m)).(configs.Server))
	}

	container.GetCoreApplication().Run(servers)
	container.GetCoreRouterGateway().Register(servers)
}
