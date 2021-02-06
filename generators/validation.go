package generators

import (
	"fmt"
	"os"
	"strings"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
)

type Validation struct {
}

func (g *Validation) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	validationTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/validation.tpl", workDir, templatePath))
	validationPath := fmt.Sprintf("%s/validations", modulePath)
	os.MkdirAll(validationPath, 0755)

	validationFile, err := os.Create(fmt.Sprintf("%s/%s.go", validationPath, strings.ToLower(template.Module)))
	if err != nil {
		panic(err)
	}

	validationTemplate.Execute(validationFile, template)
}
