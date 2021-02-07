package {{.ModulePluralLowercase}}

import (
    "context"

	configs "{{.PackageName}}/configs"
	grpcs "{{.PackageName}}/protos/builds"
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	Env      *configs.Env
	Module   *Module
	Database *gorm.DB
}

func (s *Server) RegisterGRpc(gs *grpc.Server) {
	grpcs.Register{{.ModulePlural}}Server(gs, s.Module)
}

func (s *Server) GRpcHandler(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	return grpcs.Register{{.ModulePlural}}Handler(context, server, client)
}

func (s *Server) RegisterAutoMigrate() {
	if s.Env.DbAutoMigrate {
		s.Database.AutoMigrate(&models.{{.Module}}{})
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
