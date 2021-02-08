package routes

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	configs "github.com/crowdeco/skeleton/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type MuxRouter struct {
}

func (m *MuxRouter) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	server.HandleFunc("/api/docs/", func(w http.ResponseWriter, r *http.Request) {
		regex := regexp.MustCompile("/api/docs/")
		http.ServeFile(w, r, regex.ReplaceAllString(r.URL.Path, "swaggers/"))
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

func (m *MuxRouter) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
