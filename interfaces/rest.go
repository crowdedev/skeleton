package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
	handlers "github.com/crowdeco/skeleton/handlers"
	middlewares "github.com/crowdeco/skeleton/middlewares"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
)

type (
	Rest struct{}
)

func NewRest() configs.Application {
	return &Rest{}
}

func (g *Rest) Run() {
	log.Printf("Starting REST Server on :%d", configs.Env.HtppPort)

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	endpoint := fmt.Sprintf("0.0.0.0:%d", configs.Env.RpcPort)
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/apidocs/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/apidocs/")
		p = path.Join("swagger", p)
		http.ServeFile(w, r, p)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "OK")
	})

	gw, err := newGateway(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", gw)

	middleware := handlers.NewMiddleware()
	middleware.Add(middlewares.NewAuth())

	log.Printf("API Documentation is ready at: http://localhost:%d/apidocs/ui", configs.Env.HtppPort)

	http.ListenAndServe(fmt.Sprintf(":%d", configs.Env.HtppPort), middleware.Attach(mux))
}
