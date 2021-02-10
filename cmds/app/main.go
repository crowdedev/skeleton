package main

import (
	"fmt"

	"github.com/crowdeco/skeleton/configs"
	"github.com/crowdeco/skeleton/generated/dic"
	"github.com/sirupsen/logrus"
)

func main() {
	container, _ := dic.NewContainer()
	parser := container.GetCoreConfigParser()

	var servers []configs.Server
	for _, c := range parser.ParseModules() {
		servers = append(servers, container.Get(fmt.Sprintf("%s:server", c)).(configs.Server))
	}

	var listeners []configs.Listener
	for _, c := range parser.ParseListeners() {
		listeners = append(listeners, container.Get(c).(configs.Listener))
	}

	var middlewares []configs.Middleware
	for _, c := range parser.ParseMiddlewares() {
		middlewares = append(middlewares, container.Get(c).(configs.Middleware))
	}

	var extensions []logrus.Hook
	for _, c := range parser.ParseLoggers() {
		extensions = append(extensions, container.Get(c).(logrus.Hook))
	}

	var routes []configs.Route
	for _, c := range parser.ParseRoutes() {
		routes = append(routes, container.Get(c).(configs.Route))
	}

	container.GetCoreRouterMux().Register(routes)
	container.GetCoreLoggerExtension().Register(extensions)
	container.GetCoreHandlerMiddleware().Register(middlewares)
	container.GetCoreEventDispatcher().Register(listeners)
	container.GetCoreRouterGateway().Register(servers)
	container.GetCoreApplication().Run(servers)
}
