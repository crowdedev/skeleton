package routes

import (
	"context"
	"net/http"

	grpcs "github.com/crowdeco/skeleton/protos/builds"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRpcGateway struct {
}

func (g *GRpcGateway) Handle(ctx context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	var handlers []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

	mux := runtime.NewServeMux()
	handlers = append(handlers, grpcs.RegisterTodosHandler)

	for _, handler := range handlers {
		handler(ctx, mux, client)
	}

	server.Handle("/", mux)

	return server
}
