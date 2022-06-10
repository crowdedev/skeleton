package dics

import (
	"github.com/KejawenLab/skeleton/todos"
	"github.com/KejawenLab/skeleton/todos/models"
	"github.com/KejawenLab/skeleton/todos/validations"
	"github.com/sarulabs/dingo/v4"
)

var Todo = []dingo.Def{
	{
		Name:  "module:todo:model",
		Build: (*models.Todo)(nil),
        Params: dingo.Params{
			"GormModel": dingo.Service("bima:model"),
		},
	},
	{
		Name:  "module:todo:validation",
		Build: (*validations.Todo)(nil),
	},
	{
		Name:  "module:todo",
		Build: (*todos.Module)(nil),
		Params: dingo.Params{
			"Module":    dingo.Service("bima:module"),
			"Validator": dingo.Service("module:todo:validation"),
		},
	},
	{
		Name:  "module:todo:server",
		Build: (*todos.Server)(nil),
		Params: dingo.Params{
			"Server": dingo.Service("bima:server"),
			"Module": dingo.Service("module:todo"),
		},
	},
}
