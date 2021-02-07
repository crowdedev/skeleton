package routes

import (
	"context"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRpcGateway struct {
	Servers []configs.Server
}

func (g *GRpcGateway) Register(servers []configs.Server) {
	g.Servers = servers
}

func (g *GRpcGateway) Handle(ctx context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	mux := runtime.NewServeMux()

	for _, handler := range g.Servers {
		handler.GRpcHandler(ctx, mux, client)
	}

	server.Handle("/", mux)

	return server
}

func (a *GRpcGateway) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
