package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	"github.com/sarulabs/dingo/v4"
	"google.golang.org/grpc"
)

var Interface = []dingo.Def{
	{
		Name: "core:interface:database",
		Build: func(
			todo configs.Server,
		) (*interfaces.Database, error) {
			database := interfaces.Database{
				Servers: []configs.Server{
					todo,
				},
			}

			return &database, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("module:todo:server"),
		},
	},
	{
		Name: "core:interface:grpc",
		Build: func(
			env *configs.Env,
			todo configs.Server,
			server *grpc.Server,
			dispatcher *events.Dispatcher,
		) (*interfaces.GRpc, error) {
			grpc := interfaces.GRpc{
				Env:        env,
				GRpc:       server,
				Dispatcher: dispatcher,
			}

			grpc.Register([]configs.Server{
				todo,
			})

			return &grpc, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:config:env"),
			"1": dingo.Service("module:todo:server"),
		},
	},
	{
		Name: "core:interface:queue",
		Build: func(
			todo configs.Server,
		) (*interfaces.Queue, error) {
			queue := interfaces.Queue{
				Servers: []configs.Server{
					todo,
				},
			}

			return &queue, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("module:todo:server"),
		},
	},
}
