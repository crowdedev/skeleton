package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	routes "github.com/crowdeco/skeleton/routes"
	"github.com/sarulabs/dingo/v4"
)

var Router = []dingo.Def{
	{
		Name: "core:router:gateway",
		Build: func(
			bank configs.Server,
		) (*routes.GRpcGateway, error) {
			return &routes.GRpcGateway{[]configs.Server{
				bank,
			}}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("module:bank:server"),
		},
	},
}
