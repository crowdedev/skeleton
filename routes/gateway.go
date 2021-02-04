package routes

import (
	"context"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type RegisterHandler = func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type gRpcGateway struct {
	context  context.Context
	client   *grpc.ClientConn
	handlers []RegisterHandler
}

func NewGRpcGateway(context context.Context, client *grpc.ClientConn, handlers []RegisterHandler) configs.Router {
	return &gRpcGateway{
		context:  context,
		client:   client,
		handlers: handlers,
	}
}

func (g *gRpcGateway) Handle(server *http.ServeMux) *http.ServeMux {
	mux := runtime.NewServeMux()

	for _, handler := range g.handlers {
		handler(g.context, mux, g.client)
	}

	server.Handle("/", mux)

	return server
}
