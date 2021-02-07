package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
	"gopkg.in/yaml.v2"
)

type Module struct {
	Config *configs.Config
}

func (g *Module) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	moduleTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/module.tpl", workDir, templatePath))
	moduleFile, err := os.Create(fmt.Sprintf("%s/module.go", modulePath))
	if err != nil {
		panic(err)
	}

	g.Config.Parse()
	g.Config.Modules = append(g.Config.Modules, fmt.Sprintf("module:%s", template.ModuleLowercase))

	modules, err := yaml.Marshal(g.Config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(configs.MODULES_FILE, modules, 0644)
	if err != nil {
		panic(err)
	}

	moduleTemplate.Execute(moduleFile, template)
}
