package main

import (
	configs "github.com/crowdeco/skeleton/configs"
)

func init() {
	configs.LoadConfigs()
	configs.Env.ServiceName = "skeleton"
	configs.Env.Version = "v2.0@dev"
}

func main() {
	//
}
