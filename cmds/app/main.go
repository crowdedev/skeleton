package main

import (
	"fmt"
	"time"

	configs "github.com/crowdeco/bima/configs"
	dic "github.com/crowdeco/skeleton/generated/dic"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := container.GetBimaUtilCli()

	util.Println("✍  Engine Checking and Configuring...")
	time.Sleep(1 * time.Second)

	var servers []configs.Server
	for _, c := range container.GetBimaConfigParserModule().Parse() {
		servers = append(servers, container.Get(fmt.Sprintf("%s:server", c)).(configs.Server))
	}

	var listeners []configs.Listener
	for _, c := range container.GetBimaConfigParserListener().Parse() {
		listeners = append(listeners, container.Get(c).(configs.Listener))
	}

	var middlewares []configs.Middleware
	for _, c := range container.GetBimaConfigParserMiddleware().Parse() {
		middlewares = append(middlewares, container.Get(c).(configs.Middleware))
	}

	var extensions []logrus.Hook
	for _, c := range container.GetBimaConfigParserLogger().Parse() {
		extensions = append(extensions, container.Get(c).(logrus.Hook))
	}

	var routes []configs.Route
	for _, c := range container.GetBimaConfigParserRoute().Parse() {
		routes = append(routes, container.Get(c).(configs.Route))
	}

	util.Println("🔥 Engine Start...")
	time.Sleep(3 * time.Second)

	container.GetBimaRouterMux().Register(routes)
	container.GetBimaLoggerExtension().Register(extensions)
	container.GetBimaHandlerMiddleware().Register(middlewares)
	container.GetBimaEventDispatcher().Register(listeners)
	container.GetBimaRouterGateway().Register(servers)

	util.Println("🚀 Taking Off...")
	time.Sleep(1 * time.Second)

	util.Println("🎧 🎧 🎧 Enjoy The Flight 🎧 🎧 🎧")

	container.GetBimaApplication().Run(servers)
}
