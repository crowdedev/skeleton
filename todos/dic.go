package todos

import (
	"github.com/KejawenLab/bima/v4"
    "github.com/sarulabs/dingo/v4"
)

var Dic = []dingo.Def{
	{
		Name:  "module:todo:model",
        Scope: bima.Application,
		Build: (*Todo)(nil),
        Params: dingo.Params{
			"GormModel": dingo.Service("bima:model"),
		},
	},
	{
		Name:  "module:todo",
        Scope: bima.Application,
		Build: (*Module)(nil),
		Params: dingo.Params{
            "Model":  dingo.Service("module:todo:model"),
			"Module": dingo.Service("bima:module"),
		},
	},
	{
		Name:  "module:todo:server",
        Scope: bima.Application,
		Build: (*Server)(nil),
		Params: dingo.Params{
			"Server": dingo.Service("bima:server"),
			"Module": dingo.Service("module:todo"),
		},
	},
}
