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

type Rest struct {
	Env        *configs.Env
	Middleware *handlers.Middleware
	Router     *handlers.Router
	Server     *http.ServeMux
	Context    context.Context
}

func (r *Rest) Run(servers []configs.Server) {
	log.Printf("Starting REST Server on :%d", r.Env.HtppPort)

	ctx, cancel := context.WithCancel(r.Context)
	defer cancel()

	endpoint := fmt.Sprintf("0.0.0.0:%d", r.Env.RpcPort)
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

	log.Println("API Documentation is ready at /api/docs")

	http.ListenAndServe(fmt.Sprintf(":%d", r.Env.HtppPort), r.Middleware.Attach(r.Router.Handle(ctx, r.Server, conn)))
}

func (r *Rest) IsBackground() bool {
	return false
}

func (r *Rest) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
