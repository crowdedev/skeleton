package main

import (
	configs "github.com/crowdeco/todo-service/configs"
	consoles "github.com/crowdeco/todo-service/consoles"
)

func init() {
	configs.LoadConfigs()
	configs.Env.ServiceName = "todo-service"
	configs.Env.Version = "v2.0@dev"
}

func main() {
    //TODO
}
