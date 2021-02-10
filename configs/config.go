package configs

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const MODULES_FILE = "modules.yaml"
const LISTENERS_FILE = "listeners.yaml"
const LOGGERS_FILE = "loggers.yaml"
const MIDDLEWARES_FILE = "middlewares.yaml"

type Config struct {
	Modules     []string `yaml:"modules"`
	Listeners   []string `yaml:"listeners"`
	Loggers     []string `yaml:"loggers"`
	Middlewares []string `yaml:"middlewares"`
}

func (c *Config) ParseModules() []string {
	c.parse(MODULES_FILE)

	return c.Modules
}

func (c *Config) ParseListeners() []string {
	c.parse(LISTENERS_FILE)

	return c.Listeners
}

func (c *Config) ParseLoggers() []string {
	c.parse(LOGGERS_FILE)

	return c.Loggers
}

func (c *Config) ParseMiddlewares() []string {
	c.parse(MIDDLEWARES_FILE)

	return c.Middlewares
}

func (c *Config) parse(file string) {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, file))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &c)
	if err != nil {
		panic(err)
	}
}
