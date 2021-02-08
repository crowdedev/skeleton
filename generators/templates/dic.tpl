package modules

import (
	configs "{{.PackageName}}/configs"
	{{.ModulePluralLowercase}} "{{.PackageName}}/{{.ModulePluralLowercase}}"
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	services "{{.PackageName}}/{{.ModulePluralLowercase}}/services"
	validations "{{.PackageName}}/{{.ModulePluralLowercase}}/validations"
	"github.com/sarulabs/dingo/v4"
	"gorm.io/gorm"
)

var {{.Module}} = []dingo.Def{
	{
		Name:  "module:{{.ModuleLowercase}}:model",
		Build: (*models.{{.Module}})(nil),
	},
	{
		Name:  "module:{{.ModuleLowercase}}:validation",
		Build: (*validations.{{.Module}})(nil),
	},
	{
		Name:  "module:{{.ModuleLowercase}}",
		Build: (*{{.ModulePluralLowercase}}.Module)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
			"Handler":       dingo.Service("core:handler:handler"),
			"Logger":        dingo.Service("core:handler:logger"),
			"Messenger":     dingo.Service("core:handler:messager"),
			"Validator":     dingo.Service("module:{{.ModuleLowercase}}:validation"),
			"Cache":         dingo.Service("core:cache:memory"),
			"Paginator":     dingo.Service("core:pagination:paginator"),
		},
	},
	{
		Name:  "module:{{.ModuleLowercase}}:server",
		Build: (*{{.ModulePluralLowercase}}.Server)(nil),
		Params: dingo.Params{
			"Env":      dingo.Service("core:config:env"),
			"Module":   dingo.Service("module:{{.ModuleLowercase}}"),
			"Database": dingo.Service("core:connection:database"),
		},
	},
}
