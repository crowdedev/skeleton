package generators

import (
	"fmt"
	"os"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
)

type Server struct {
}

func (g *Server) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	serverTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/server.tpl", workDir, templatePath))
	serverFile, err := os.Create(fmt.Sprintf("%s/server.go", modulePath))
	if err != nil {
		panic(err)
	}

	serverTemplate.Execute(serverFile, template)
}
