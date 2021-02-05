package configs

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type (
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
	}

	Module interface {
		Consume()
	}

	Server interface {
		RegisterGRpc(server *grpc.Server)
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
