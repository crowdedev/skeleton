package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	generators "github.com/crowdeco/skeleton/generators"
	"github.com/gertd/go-pluralize"
	"github.com/sarulabs/dingo/v4"
)

var Generator = []dingo.Def{
	{
		Name: "core:module:generator",
		Build: func(
			model configs.Generator,
			module configs.Generator,
			proto configs.Generator,
			server configs.Generator,
			service configs.Generator,
			validation configs.Generator,
			pluralizer *pluralize.Client,
			template *configs.Template,
		) (*generators.Factory, error) {
			return &generators.Factory{
				Pluralizer: pluralizer,
				Template:   template,
				Generators: []configs.Generator{
					model,
					module,
					proto,
					server,
					service,
					validation,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:generator:model"),
			"1": dingo.Service("core:generator:module"),
			"2": dingo.Service("core:generator:proto"),
			"3": dingo.Service("core:generator:server"),
			"4": dingo.Service("core:generator:service"),
			"5": dingo.Service("core:generator:validation"),
		},
	},
}
