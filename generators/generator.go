package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	configs "github.com/crowdeco/skeleton/configs"
	"github.com/gertd/go-pluralize"
	"golang.org/x/mod/modfile"
)

const TEMPLATE_PATH = "generators/templates"

type Generator struct {
	Pluralizer *pluralize.Client
	Template   *configs.Template
}

func (g *Generator) Generate(module string, templatePath string) {
	workDir, _ := os.Getwd()
	packageName := g.GetPackageName(workDir)
	moduleName := strings.Title(module)
	moduleLowercase := strings.ToLower(module)
	modulePlural := g.Pluralizer.Plural(moduleName)
	modulePluralLowercase := strings.ToLower(modulePlural)
	modulePath := fmt.Sprintf("%s/%s", workDir, modulePluralLowercase)

	g.Template.PackageName = packageName
	g.Template.Module = moduleName
	g.Template.ModulePlural = modulePlural
	g.Template.ModulePluralLowercase = modulePluralLowercase

	serviceTemplate, _ := template.ParseFiles(fmt.Sprintf("%s/%s/service.tpl", workDir, templatePath))
	modelTemplate, _ := template.ParseFiles(fmt.Sprintf("%s/%s/model.tpl", workDir, templatePath))
	validationTemplate, _ := template.ParseFiles(fmt.Sprintf("%s/%s/validation.tpl", workDir, templatePath))
	moduleTemplate, _ := template.ParseFiles(fmt.Sprintf("%s/%s/module.tpl", workDir, templatePath))
	serverTemplate, _ := template.ParseFiles(fmt.Sprintf("%s/%s/server.tpl", workDir, templatePath))
	protoTemplate, _ := template.ParseFiles(fmt.Sprintf("%s/%s/proto.tpl", workDir, templatePath))

	serviceFile, modelFile, validationFile, moduleFile, serverFile, protoFile := g.getFiles(workDir, modulePath, moduleLowercase)

	serviceTemplate.Execute(serviceFile, g.Template)
	modelTemplate.Execute(modelFile, g.Template)
	validationTemplate.Execute(validationFile, g.Template)
	moduleTemplate.Execute(moduleFile, g.Template)
	serverTemplate.Execute(serverFile, g.Template)
	protoTemplate.Execute(protoFile, g.Template)
}

func (g *Generator) GetDefaultTemplatePath() string {
	return TEMPLATE_PATH
}

func (g *Generator) GetPackageName(workDir string) string {
	goModBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	return modfile.ModulePath(goModBytes)
}

func (g *Generator) createModuleDirectory(modulePath string) (string, string, string) {
	servicePath := fmt.Sprintf("%s/services", modulePath)
	modelPath := fmt.Sprintf("%s/models", modulePath)
	validationPath := fmt.Sprintf("%s/validations", modulePath)

	os.MkdirAll(servicePath, 0755)
	os.MkdirAll(modelPath, 0755)
	os.MkdirAll(validationPath, 0755)

	return servicePath, modelPath, validationPath
}

func (g *Generator) getFiles(workDir string, modulePath string, moduleLowercase string) (*os.File, *os.File, *os.File, *os.File, *os.File, *os.File) {
	servicePath, modelPath, validationPath := g.createModuleDirectory(modulePath)

	serviceFile, err := os.Create(fmt.Sprintf("%s/%s.go", servicePath, moduleLowercase))
	if err != nil {
		panic(err)
	}

	modelFile, err := os.Create(fmt.Sprintf("%s/%s.go", modelPath, moduleLowercase))
	if err != nil {
		panic(err)
	}

	validationFile, err := os.Create(fmt.Sprintf("%s/%s.go", validationPath, moduleLowercase))
	if err != nil {
		panic(err)
	}

	moduleFile, err := os.Create(fmt.Sprintf("%s/module.go", modulePath))
	if err != nil {
		panic(err)
	}

	serverFile, err := os.Create(fmt.Sprintf("%s/server.go", modulePath))
	if err != nil {
		panic(err)
	}

	protoFile, err := os.Create(fmt.Sprintf("%s/protos/%s.proto", workDir, moduleLowercase))
	if err != nil {
		panic(err)
	}

	return serviceFile, modelFile, validationFile, moduleFile, serverFile, protoFile
}
