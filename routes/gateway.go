package routes

import (
	"context"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type gRpcGateway struct {
	context context.Context
	client  *grpc.ClientConn
}

func NewGRpcGateway(context context.Context, client *grpc.ClientConn) configs.Router {
	return &gRpcGateway{
		context: context,
		client:  client,
	}
}

func (g *gRpcGateway) Handle(server *http.ServeMux) *http.ServeMux {
	var handlers []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

	mux := runtime.NewServeMux()
	handlers = append(handlers, grpcs.RegisterTodosHandler)

	for _, handler := range handlers {
		handler(g.context, mux, g.client)
	}

	server.Handle("/", mux)

	return server
}
