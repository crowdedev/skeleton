package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	handlers "github.com/crowdeco/skeleton/handlers"
	middlewares "github.com/crowdeco/skeleton/middlewares"
	routes "github.com/crowdeco/skeleton/routes"
	"google.golang.org/grpc"
)

type (
	Client struct {
		Grpc    *grpc.ClientConn
		Context context.Context
	}

	Rest struct {
		Client *Client
	}
)

func (g *Rest) Run() {
	log.Printf("Starting REST Server on :%d", configs.Env.HtppPort)

	mux := http.NewServeMux()

	router := handlers.NewRouter()
	router.Add(routes.NewMuxRouter(g.Client.Grpc))
	router.Add(routes.NewGRpcGateway(g.Client.Context, g.Client.Grpc))

	middleware := handlers.NewMiddleware()
	middleware.Add(middlewares.NewAuth())

	log.Println("API Documentation is ready at /api/docs/ui")

	http.ListenAndServe(fmt.Sprintf(":%d", configs.Env.HtppPort), middleware.Attach(router.Handle(mux)))
}
