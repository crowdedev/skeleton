package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	handlers "github.com/crowdeco/skeleton/handlers"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type (
	Rest struct {
		Middleware *handlers.Middleware
		Router     *handlers.Router
		Server     *http.ServeMux
		Context    context.Context
	}
)

func (g *Rest) Run() {
	log.Printf("Starting REST Server on :%d", configs.Env.HtppPort)

	ctx, cancel := context.WithCancel(g.Context)
	defer cancel()

	endpoint := fmt.Sprintf("0.0.0.0:%d", configs.Env.RpcPort)
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
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

	log.Println("API Documentation is ready at /api/docs/ui")

	http.ListenAndServe(fmt.Sprintf(":%d", configs.Env.HtppPort), g.Middleware.Attach(g.Router.Handle(ctx, g.Server, conn)))
}
