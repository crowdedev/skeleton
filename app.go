package skeleton

import (
	"fmt"
	"os"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/parsers"
	"github.com/KejawenLab/skeleton/generated/dic"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Run() {
	workDir, _ := os.Getwd()
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := color.New(color.FgCyan, color.Bold)
	env := container.GetBimaConfigEnv()

	if env.Debug {
		util.Println("✓  Engine Checking and Configuring...")
		time.Sleep(500 * time.Millisecond)
	}

	var servers []configs.Server
	for _, c := range parsers.ParseModule(workDir) {
		servers = append(servers, container.Get(fmt.Sprintf("%s:server", c)).(configs.Server))
	}

	var listeners []configs.Listener
	for _, c := range parsers.ParseListener(workDir) {
		listeners = append(listeners, container.Get(fmt.Sprintf("bima:listener:%s", c)).(configs.Listener))
	}

	var middlewares []configs.Middleware
	for _, c := range parsers.ParseMiddleware(workDir) {
		middlewares = append(middlewares, container.Get(fmt.Sprintf("bima:middleware:%s", c)).(configs.Middleware))
	}

	var extensions []logrus.Hook
	for _, c := range parsers.ParseLogger(workDir) {
		extensions = append(extensions, container.Get(fmt.Sprintf("bima:logger:extension:%s", c)).(logrus.Hook))
	}

	var routes []configs.Route
	for _, c := range parsers.ParseRoute(workDir) {
		routes = append(routes, container.Get(fmt.Sprintf("bima:route:%s", c)).(configs.Route))
	}

	if env.Debug {
		util.Printf("✓ ")
		fmt.Printf("Total pessanger %d\n", len(servers))
		util.Println("⌛ Starting Engine...")
		time.Sleep(100 * time.Millisecond)
		util.Println("🔥🔥🔥🔥🔥🔥🔥🔥🔥")
		time.Sleep(100 * time.Millisecond)
		util.Println("🔥🔥🔥🔥🔥🔥🔥🔥🔥")
		time.Sleep(100 * time.Millisecond)
		util.Println("🔥🔥🔥🔥🔥🔥🔥🔥🔥")
		time.Sleep(100 * time.Millisecond)
		util.Println("🔥 Engine Ready...")
		time.Sleep(1 * time.Second)
	}

	container.GetBimaRouterMux().Register(routes)
	container.GetBimaLoggerExtension().Register(extensions)
	container.GetBimaHandlerMiddleware().Register(middlewares)
	container.GetBimaEventDispatcher().Register(listeners)
	container.GetBimaRouterGateway().Register(servers)

	if env.Debug {
		util.Println("🚀 Taking Off...")
		time.Sleep(300 * time.Millisecond)

		util.Println("🎧 🎧 🎧 Enjoy The Flight 🎧 🎧 🎧")
	}

	container.GetBimaApplication().Run(servers)
}
