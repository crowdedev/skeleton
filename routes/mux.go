package routes

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type muxRouter struct {
	client *grpc.ClientConn
}

func NewMuxRouter(client *grpc.ClientConn) configs.Router {
	return &muxRouter{
		client: client,
	}
}

func (g *muxRouter) Handle(server *http.ServeMux) *http.ServeMux {
	server.HandleFunc("/api/docs/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/api/docs/")
		p = path.Join("swagger", p)
		http.ServeFile(w, r, p)
	})

	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		s := g.client.GetState()

		if s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("gRPC server is %s", s), http.StatusBadGateway)

			return
		}

		fmt.Fprintln(w, "OK")
	})

	return server
}
