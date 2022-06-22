package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/KejawenLab/skeleton"
)

func main() {
	args := os.Args[1:]
	command := args[0]
	option := ""
	module := ""
	if len(args) > 1 {
		option = args[1]
	}
	if len(args) > 2 {
		module = args[2]
	}

	switch command {
	case "dump":
		_, err := exec.Command("go", "run", "dumper/main.go").Output()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	case "run":
		skeleton.Application(command).Run()
	case "module":
		skeleton.Module(option).Run(module)
	}
}
