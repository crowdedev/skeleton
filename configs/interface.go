package configs

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type (
	Driver interface {
		Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB
	}

	Generator interface {
		Generate(template *Template, modulePath string, workDir string, templatePath string)
	}

	Listener interface {
		Handle(event interface{})
		Listen() string
	}

	Model interface {
		TableName() string
		SetCreatedBy(user *User)
		SetUpdatedBy(user *User)
		SetDeletedBy(user *User)
		IsSoftDelete() bool
	}

	Service interface {
		Name() string
		Create(value interface{}, id string) error
		Update(value interface{}, id string) error
		Bind(value interface{}, id string) error
		Delete(value interface{}, id string) error
		All(value interface{}) error
	}

	Module interface {
		Consume()
	}

	Server interface {
		RegisterGRpc(server *grpc.Server)
		GRpcHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
		RegisterAutoMigrate()
		RegisterQueueConsumer()
	}

	Router interface {
		Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux
	}

	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
	}

	Application interface {
		Run()
	}
)
