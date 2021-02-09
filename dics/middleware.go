package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	handlers "github.com/crowdeco/skeleton/handlers"
	"github.com/sarulabs/dingo/v4"
)

var Middleware = []dingo.Def{
	{
		Name: "core:handler:middleware",
		Build: func(
			auth configs.Middleware,
		) (*handlers.Middleware, error) {
			return &handlers.Middleware{
				Middlewares: []configs.Middleware{
					auth,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:middleware:auth"),
		},
	},
}
