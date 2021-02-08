package configs

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const HIGEST_PRIORITY = 255

const LOWEST_PRIORITY = -255

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
		Priority() int
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
		OverrideData(value interface{})
		Create(value interface{}) error
		Update(value interface{}, id string) error
		Bind(value interface{}, id string) error
		All(value interface{}) error
		Delete(value interface{}, id string) error
	}

	Module interface {
		Consume()
		Populete()
	}

	Server interface {
		RegisterGRpc(server *grpc.Server)
		GRpcHandler(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error
		RegisterAutoMigrate()
		RegisterQueueConsumer()
		RepopulateData()
	}

	Router interface {
		Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux
		Priority() int
	}

	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
		Priority() int
	}

	Application interface {
		Run(servers []Server)
		IsBackground() bool
		Priority() int
	}
)
