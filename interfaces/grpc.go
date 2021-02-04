package interfaces

import (
	"fmt"
	"log"
	"net"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	parents "github.com/crowdeco/skeleton/parents"
	todos "github.com/crowdeco/skeleton/todos"
	grpc "google.golang.org/grpc"
)

type gRpc struct {
	dispatcher *events.Dispatcher
}

func NewGRpc(dispatcher *events.Dispatcher) configs.Application {
	return &gRpc{dispatcher: dispatcher}
}

func (g *gRpc) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Env.RpcPort))
	if err != nil {
		log.Fatalf("Port %d is not available. %v", configs.Env.RpcPort, err)
	}

	app := grpc.NewServer()
	parents.NewServer(g.dispatcher).RegisterGRpc(app)
	todos.NewServer(g.dispatcher).RegisterGRpc(app)

	log.Printf("Starting gRPC Server on :%d", configs.Env.RpcPort)

	app.Serve(l)
}
