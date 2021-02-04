package todos

import (
	configs "github.com/crowdeco/skeleton/configs"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	Env      *configs.Env
	Module   *TodoModule
	Database *gorm.DB
}

func (s *Server) RegisterGRpc(gs *grpc.Server) {
	grpcs.RegisterTodosServer(gs, s.Module)
}

func (s *Server) RegisterAutoMigrate() {
	if s.Env.DbAutoMigrate {
		s.Database.AutoMigrate(&models.Todo{})
	}
}

func (s *Server) RegisterQueueConsumer() {
	s.Module.Consume()
}
