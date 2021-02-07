package main

import (
	dic "github.com/crowdeco/skeleton/generated/dic"
)

func main() {
	container, _ := dic.NewContainer()

	application, err := container.SafeGetCoreApplication()
	if err != nil {
		panic(err)
	}

	application.Run()
}
