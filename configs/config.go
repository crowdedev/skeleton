package configs

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const MODULES_FILE = "modules.yaml"

type Config struct {
	Modules []string `yaml:"modules"`
}

func (c *Config) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, MODULES_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &c)
	if err != nil {
		panic(err)
	}

	return c.Modules
}
