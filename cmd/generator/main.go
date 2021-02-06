package main

import (
	"fmt"
	"os"

	"github.com/crowdeco/skeleton/generated/dic"
)

func main() {
	container, _ := dic.NewContainer()
	generator, err := container.SafeGetCoreModuleGenerator()
	if err != nil {
		fmt.Errorf("Error Generator: %s", err.Error())
		return
	}

	templatePath := generator.GetDefaultTemplatePath()
	if len(os.Args) == 2 {
		templatePath = os.Args[1]
	}

	generator.Generate("todo", templatePath)
}
