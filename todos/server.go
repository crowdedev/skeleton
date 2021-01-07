package todos

import (
	"context"
	"fmt"

	configs "github.com/crowdeco/todo-service/configs"
	grpcs "github.com/crowdeco/todo-service/protos/builds"
	models "github.com/crowdeco/todo-service/todos/models"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterGrpcServer(s *grpc.Server) {
	grpcs.RegisterTodosServer(s, NewTodo())
}

func RegisterRestServer(c context.Context, r *runtime.ServeMux) {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := grpcs.RegisterTodosHandlerFromEndpoint(c, r, fmt.Sprintf("0.0.0.0:%d", configs.Env.RpcPort), opts)
	if err != nil {
		panic(err)
	}
}

func RegisterAutoMigration() {
	if configs.Env.DbAutoMigrate {
		configs.Database.AutoMigrate(&models.Todo{})
	}
}
