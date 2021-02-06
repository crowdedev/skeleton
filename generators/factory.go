package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
	"github.com/gertd/go-pluralize"
	"golang.org/x/mod/modfile"
)

const TEMPLATE_PATH = "generators/templates"

type Factory struct {
	Pluralizer *pluralize.Client
	Template   *configs.Template
	Generators []configs.Generator
}

func (f *Factory) Generate(module string, templatePath string) {
	workDir, _ := os.Getwd()
	packageName := f.GetPackageName(workDir)
	moduleName := strings.Title(module)
	modulePlural := f.Pluralizer.Plural(moduleName)
	modulePluralLowercase := strings.ToLower(modulePlural)
	modulePath := fmt.Sprintf("%s/%s", workDir, modulePluralLowercase)

	f.Template.PackageName = packageName
	f.Template.Module = moduleName
	f.Template.ModulePlural = modulePlural
	f.Template.ModulePluralLowercase = modulePluralLowercase

	for _, generator := range f.Generators {
		generator.Generate(f.Template, modulePath, workDir, templatePath)
	}
}

func (f *Factory) GetDefaultTemplatePath() string {
	return TEMPLATE_PATH
}

func (f *Factory) GetPackageName(workDir string) string {
	goModBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	return modfile.ModulePath(goModBytes)
}
