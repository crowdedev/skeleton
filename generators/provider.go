package generators

import (
	"fmt"
	"io/ioutil"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
)

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
		if strings.Contains(v, "//@modules:import") {
			importIdx = k
			skipImport = false
			continue
		}

		if strings.Contains(v, "//@modules:register") {
			moduleIdx = k
			break
		}
	}

	if !skipImport {
		contents[importIdx] = fmt.Sprintf("    modules %q", "github.com/crowdeco/skeleton/dics/modules")
	}

	contents[moduleIdx] = fmt.Sprintf(`
    @module:%sif err := p.AddDefSlice(modules.%s); err != nil {return err}
    //@modules:register`, template.ModuleLowercase, template.Module)

	body := strings.Join(contents, "\n")

	ioutil.WriteFile(path, []byte(body), 0644)
}
