package main

import (
	configs "github.com/crowdeco/skeleton/configs"
	consoles "github.com/crowdeco/skeleton/consoles"
)

func init() {
	configs.LoadConfigs()
	configs.Env.ServiceName = "skeleton"
	configs.Env.Version = "v2.0@dev"
}

func main() {
    //TODO
}
