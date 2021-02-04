package handlers

import (
	"context"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	"google.golang.org/grpc"
)

type Router struct {
	Routes []configs.Router
}

func (r *Router) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	for _, route := range r.Routes {
		route.Handle(context, server, client)
	}

	return server
}
