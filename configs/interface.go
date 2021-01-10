package configs

import (
	"context"
	"net/http"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	Model interface {
		TableName() string
		Identifier() string
		SetIdentifier(id string)
		SetCreatedBy(user *User)
		SetUpdatedBy(user *User)
		SetDeletedBy(user *User)
		IsSoftDelete() bool
	}

	Service interface {
		Model() Model
		Create() Model
		Update() Model
		Bind() Model
		Delete()
	}

	Server interface {
		RegisterRest(context context.Context, runtime *runtime.ServeMux)
		RegisterGRpc(server *grpc.Server)
		RegisterAutoMigrate()
		RegisterQueueConsumer()
	}

	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
	}

	Application interface {
		Run()
	}
)
