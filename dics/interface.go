package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	"github.com/sarulabs/dingo/v4"
)

var Interface = []dingo.Def{
	{
		Name: "core:application",
		Build: func(
			database configs.Application,
			elasticsearch configs.Application,
			grpc configs.Application,
			queue configs.Application,
			rest configs.Application,
		) (*interfaces.Application, error) {
			return &interfaces.Application{
				Applications: []configs.Application{database, elasticsearch, grpc, queue, rest},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:interface:database"),
			"1": dingo.Service("core:interface:elasticsearch"),
			"2": dingo.Service("core:interface:grpc"),
			"3": dingo.Service("core:interface:queue"),
			"4": dingo.Service("core:interface:rest"),
		},
	},
}
