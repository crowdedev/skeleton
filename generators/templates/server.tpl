package {{.ModulePluralLowercase}}

import (
	configs "{{.PackageName}}/configs"
	grpcs "{{.PackageName}}/protos/builds"
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
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

func (s *Server) RegisterAutoMigrate() {
	if s.Env.DbAutoMigrate {
		s.Database.AutoMigrate(&models.{{.Module}}{})
	}
}

func (s *Server) RegisterQueueConsumer() {
	s.Module.Consume()
}
