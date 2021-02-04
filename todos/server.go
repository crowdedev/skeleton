package todos

import (
	configs "github.com/crowdeco/skeleton/configs"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	"google.golang.org/grpc"
)

type Server struct {
	Module *TodoModule
}

func (s *Server) RegisterGRpc(gs *grpc.Server) {
	grpcs.RegisterTodosServer(gs, s.Module)
}

func (s *Server) RegisterAutoMigrate() {
	if configs.Env.DbAutoMigrate {
		configs.Database.AutoMigrate(&models.Todo{})
	}
}

func (s *Server) RegisterQueueConsumer() {
	s.Module.Consume()
}
