package routes

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type MuxRouter struct {
}

func (g *MuxRouter) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	server.HandleFunc("/api/docs/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/api/docs/")
		p = path.Join("swagger", p)
		http.ServeFile(w, r, p)
	})

	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		s := client.GetState()

		if s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("gRPC server is %s", s), http.StatusBadGateway)

			return
		}

		fmt.Fprintln(w, "OK")
	})

	return server
}
