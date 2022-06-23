package todos

import (
    "context"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/skeleton/protos/builds"
	"github.com/KejawenLab/skeleton/todos/models"
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
	if s.Database != nil && s.Debug {
		s.Database.AutoMigrate(&models.Todo{})
	}
}

func (s *Server) RegisterQueueConsumer() {
}

func (s *Server) RepopulateData() {
}
