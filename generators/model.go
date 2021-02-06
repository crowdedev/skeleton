package generators

import (
	"fmt"
	"os"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
)

type Model struct {
}

func (g *Model) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	modelTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/model.tpl", workDir, templatePath))
	modelPath := fmt.Sprintf("%s/models", modulePath)
	os.MkdirAll(modelPath, 0755)

	modelFile, err := os.Create(fmt.Sprintf("%s/%s.go", modelPath, template.ModuleLowercase))
	if err != nil {
		panic(err)
	}

	modelTemplate.Execute(modelFile, template)
}
