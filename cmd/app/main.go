package main

import (
	dic "github.com/crowdeco/skeleton/generated/dic"
)

func main() {
	container, _ := dic.NewContainer()
	container.GetCoreApplication().Run()
}
