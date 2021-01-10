package todos

import (
	"context"
	"fmt"

	configs "github.com/crowdeco/skeleton/configs"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
}

func NewServer() configs.Server {
	return &Server{}
}

func (s *Server) RegisterGRpc(server *grpc.Server) {
	grpcs.RegisterTodosServer(server, NewTodo())
}

func (s *Server) RegisterRest(context context.Context, runtime *runtime.ServeMux) {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := grpcs.RegisterTodosHandlerFromEndpoint(context, runtime, fmt.Sprintf("0.0.0.0:%d", configs.Env.RpcPort), opts)
	if err != nil {
		panic(err)
	}
}

func (s *Server) RegisterAutoMigrate() {
	if configs.Env.DbAutoMigrate {
		configs.Database.AutoMigrate(&models.Todo{})
	}
}

func (s *Server) RegisterQueueConsumer() {
	NewTodo().Consume()
}
