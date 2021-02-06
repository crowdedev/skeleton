package generators

import (
	"fmt"
	"os"
	"strings"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
)

type Service struct {
}

func (g *Service) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	serviceTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/service.tpl", workDir, templatePath))
	servicePath := fmt.Sprintf("%s/services", modulePath)
	os.MkdirAll(servicePath, 0755)

	serviceFile, err := os.Create(fmt.Sprintf("%s/%s.go", servicePath, strings.ToLower(template.Module)))
	if err != nil {
		panic(err)
	}

	serviceTemplate.Execute(serviceFile, template)
}
