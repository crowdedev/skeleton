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
			"Model":     dingo.Service("module:todo:model"),
			"Module":    dingo.Service("bima:module"),
			"Messenger": dingo.Service("bima:messenger"),
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
