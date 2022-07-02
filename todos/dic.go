package todos

import "github.com/sarulabs/dingo/v4"

var Dic = []dingo.Def{
	{
		Name:  "module:todo:model",
		Build: (*Todo)(nil),
        Params: dingo.Params{
			"GormModel": dingo.Service("bima:model"),
		},
	},
	{
		Name:  "module:todo",
		Build: (*Module)(nil),
		Params: dingo.Params{
			"Module": dingo.Service("bima:module"),
		},
	},
	{
		Name:  "module:todo:server",
		Build: (*Server)(nil),
		Params: dingo.Params{
			"Server": dingo.Service("bima:server"),
			"Module": dingo.Service("module:todo"),
		},
	},
}
