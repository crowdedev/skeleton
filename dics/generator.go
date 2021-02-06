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
			dic configs.Generator,
			model configs.Generator,
			module configs.Generator,
			proto configs.Generator,
			server configs.Generator,
			service configs.Generator,
			validation configs.Generator,
			env *configs.Env,
			pluralizer *pluralize.Client,
			template *configs.Template,
		) (*generators.Factory, error) {
			return &generators.Factory{
				Env:        env,
				Pluralizer: pluralizer,
				Template:   template,
				Generators: []configs.Generator{
					dic,
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
			"0": dingo.Service("core:generator:dic"),
			"1": dingo.Service("core:generator:model"),
			"2": dingo.Service("core:generator:module"),
			"3": dingo.Service("core:generator:proto"),
			"4": dingo.Service("core:generator:server"),
			"5": dingo.Service("core:generator:service"),
			"6": dingo.Service("core:generator:validation"),
		},
	},
}
