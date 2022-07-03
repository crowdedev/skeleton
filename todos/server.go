package todos

import (
	"context"

	"github.com/KejawenLab/bima/v3"
	grpcs "github.com/KejawenLab/skeleton/v3/protos/builds"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	*bima.Server
	Module *Module
}

func (s *Server) Register(gs *grpc.Server) {
	grpcs.RegisterTodosServer(gs, s.Module)
}

func (s *Server) Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	return grpcs.RegisterTodosHandler(context, server, client)
}

func (s *Server) Migrate(db *gorm.DB) {
	if s.Debug {
		db.AutoMigrate(&Todo{})
	}
}
