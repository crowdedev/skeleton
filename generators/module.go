package generators

import (
	"fmt"
	"os"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
)

type Module struct {
}

func (g *Module) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	moduleTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/module.tpl", workDir, templatePath))
	moduleFile, err := os.Create(fmt.Sprintf("%s/module.go", modulePath))
	if err != nil {
		panic(err)
	}

	moduleTemplate.Execute(moduleFile, template)
}
