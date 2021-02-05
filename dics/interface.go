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
		Build: func() (*interfaces.Database, error) {
			database := interfaces.Database{
				Servers: []configs.Server{
					// @see skeleton-todo
				},
			}

			return &database, nil
		},
	},
	{
		Name: "core:interface:grpc",
		Build: func(
			env *configs.Env,
			server *grpc.Server,
			dispatcher *events.Dispatcher,
		) (*interfaces.GRpc, error) {
			grpc := interfaces.GRpc{
				Env:        env,
				GRpc:       server,
				Dispatcher: dispatcher,
			}

			grpc.Register([]configs.Server{
				// @see skeleton-todo
			})

			return &grpc, nil
		},
	},
	{
		Name: "core:interface:queue",
		Build: func() (*interfaces.Queue, error) {
			queue := interfaces.Queue{
				Servers: []configs.Server{
					// @see skeleton-todo
				},
			}

			return &queue, nil
		},
	},
}
