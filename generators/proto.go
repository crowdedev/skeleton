package generators

import (
	"fmt"
	"os"
	"strings"
	engine "text/template"

	configs "github.com/crowdeco/skeleton/configs"
)

type Proto struct {
}

func (g *Proto) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	protoTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/proto.tpl", workDir, templatePath))
	protoFile, err := os.Create(fmt.Sprintf("%s/protos/%s.proto", workDir, strings.ToLower(template.Module)))
	if err != nil {
		panic(err)
	}

	protoTemplate.Execute(protoFile, template)
}
