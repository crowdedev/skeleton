package generators

import (
	"fmt"
	"io/ioutil"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
)

const MODULE_IMPORT = "@modules:import"
const MODULE_REGISTER = "@modules:register"

type Provider struct {
}

func (p *Provider) Generate(template *configs.Template, modulePath string, workDir string, templatePath string) {
	path := fmt.Sprintf("%s/dics/provider.go", workDir)

	file, _ := ioutil.ReadFile(path)
	contents := strings.Split(string(file), "\n")
	importIdx := 0
	moduleIdx := 0
	skipImport := true

	for k, v := range contents {
		if strings.Contains(v, MODULE_IMPORT) {
			importIdx = k
			skipImport = false
			continue
		}

		if strings.Contains(v, MODULE_REGISTER) {
			moduleIdx = k
			break
		}
	}

	if !skipImport {
		contents[importIdx] = fmt.Sprintf(`    //%s
    %s %q`, MODULE_IMPORT, template.ModuleLowercase, fmt.Sprintf("%s/%s", template.PackageName, template.ModulePluralLowercase))
	}

	contents[moduleIdx] = fmt.Sprintf(`
    /*@module:%s*/if err := p.AddDefSlice(%s.%s); err != nil {return err}
    //%s`, template.ModuleLowercase, template.ModuleLowercase, template.Module, MODULE_REGISTER)

	body := strings.Join(contents, "\n")

	ioutil.WriteFile(path, []byte(body), 0644)
}
