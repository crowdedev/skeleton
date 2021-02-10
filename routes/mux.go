package routes

import (
	"context"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	"google.golang.org/grpc"
)

type MuxRouter struct {
	Routes []configs.Route
}

func (m *MuxRouter) Register(routes []configs.Route) {
	m.Routes = routes
}

func (m *MuxRouter) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	for _, v := range m.Routes {
		v.SetClient(client)
		server.HandleFunc(v.Path(), v.Handle)
	}

	return server
}

func (m *MuxRouter) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
