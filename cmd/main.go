package main

import (
	"os"

	"github.com/KejawenLab/skeleton/v3"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, "run")
	}

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
	case "run":
		skeleton.Application(command).Run(option)
	case "module":
		command = option
		option = ""
		if len(args) > 3 {
			option = args[3]
		}

		skeleton.Module(command).Run(module, option)
	}
}
