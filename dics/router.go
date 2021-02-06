package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	routes "github.com/crowdeco/skeleton/routes"
	"github.com/sarulabs/dingo/v4"
)

var Router = []dingo.Def{
	{
		Name: "core:router:gateway",
		Build: func() (*routes.GRpcGateway, error) {
			return &routes.GRpcGateway{[]configs.Server{
				// @see skeleton-todo
			}}, nil
		},
	},
}
