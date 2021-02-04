package dics

import (
	"context"
	"fmt"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
	mongodb "github.com/weekface/mgorus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"
)

var Core = []dingo.Def{
	{
		Name: "core:event:dispatcher",
		Build: func(
			todo events.Listener,
		) (*events.Dispatcher, error) {
			return events.NewDispatcher([]events.Listener{todo}), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("module:todo:listener:search"),
		},
	},
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
			todo configs.Server,
			server *grpc.Server,
			dispatcher *events.Dispatcher,
		) (*interfaces.GRpc, error) {
			grpc := interfaces.GRpc{
				GRpc:       server,
				Dispatcher: dispatcher,
			}

			grpc.Register([]configs.Server{
				todo,
			})

			return &grpc, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("module:todo:server"),
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
	{
		Name:  "core:interface:rest",
		Build: (*interfaces.Rest)(nil),
		Params: dingo.Params{
			"Client": dingo.Service("core:grpc:client"),
		},
	},
	{
		Name: "core:handler:logger",
		Build: func() (*handlers.Logger, error) {
			logger := logrus.New()
			logger.SetFormatter(&logrus.JSONFormatter{})

			mongodb, err := mongodb.NewHooker(fmt.Sprintf("%s:%d", configs.Env.MongoDbHost, configs.Env.MongoDbPort), configs.Env.MongoDbName, "logs")
			if err == nil {
				logger.AddHook(mongodb)
			} else {
				return nil, err
			}

			return &handlers.Logger{
				Logger: logger,
			}, nil
		},
	},
	{
		Name: "core:handler:messager",
		Build: func(logger *handlers.Logger) (*handlers.Messenger, error) {
			return handlers.NewMessenger(logger), nil
		},
	},
	{
		Name: "core:grpc:server",
		Build: func() (*grpc.Server, error) {
			return grpc.NewServer(), nil
		},
	},
	{
		Name: "core:context:background",
		Build: func() (context.Context, error) {
			return context.Background(), nil
		},
	},
	{
		Name: "core:grpc:client",
		Build: func(ctx context.Context) (*interfaces.Client, error) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			endpoint := fmt.Sprintf("0.0.0.0:%d", configs.Env.RpcPort)
			conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}

			defer func() {
				if err != nil {
					if cerr := conn.Close(); cerr != nil {
						grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
					}
					return
				}
				go func() {
					<-ctx.Done()
					if cerr := conn.Close(); cerr != nil {
						grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
					}
				}()
			}()

			return &interfaces.Client{
				Grpc:    conn,
				Context: ctx,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:context:background"),
		},
	},
	{
		Name: "core:gorm:db",
		Build: func() (*gorm.DB, error) {
			return configs.Database, nil
		},
	},
}
