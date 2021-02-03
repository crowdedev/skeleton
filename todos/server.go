package todos

import (
	configs "github.com/crowdeco/skeleton/configs"
	"github.com/crowdeco/skeleton/events"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	"google.golang.org/grpc"
)

type server struct {
	module TodoModule
}

func NewServer(dispatcher *events.Dispatcher) configs.Server {
	return &server{
		module: NewTodoModule(dispatcher),
	}
}

func (s *server) RegisterGRpc(gs *grpc.Server) {
	grpcs.RegisterTodosServer(gs, s.module)
}

func (s *server) RegisterAutoMigrate() {
	if configs.Env.DbAutoMigrate {
		configs.Database.AutoMigrate(&models.Todo{})
	}
}

func (s *server) RegisterQueueConsumer() {
	s.module.Consume()
}
