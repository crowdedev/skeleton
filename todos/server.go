package todos

import (
    "context"

	bima "github.com/crowdeco/bima"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
    *bima.Server
	Module *Module
}

func (s *Server) RegisterGRpc(gs *grpc.Server) {
	grpcs.RegisterTodosServer(gs, s.Module)
}

func (s *Server) GRpcHandler(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	return grpcs.RegisterTodosHandler(context, server, client)
}

func (s *Server) RegisterAutoMigrate() {
	if s.Env.DbAutoMigrate {
		s.Database.AutoMigrate(&models.Todo{})
	}
}

func (s *Server) RegisterQueueConsumer() {
	s.Module.Consume()
}

func (s *Server) RepopulateData() {
	if s.Env.Debug {
		s.Module.Populete()
	}
}
