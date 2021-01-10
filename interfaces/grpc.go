package interfaces

import (
	"fmt"
	"log"
	"net"

	configs "github.com/crowdeco/skeleton/configs"
	todos "github.com/crowdeco/skeleton/todos"
	grpc "google.golang.org/grpc"
)

type (
	GRpc struct{}
)

func NewGRpc() configs.Application {
	return &GRpc{}
}

func (g *GRpc) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Env.RpcPort))
	if err != nil {
		log.Fatalf("Port %d is not available. %v", configs.Env.RpcPort, err)
	}

	app := grpc.NewServer()
	todos.NewServer().RegisterGRpc(app)

	log.Printf("Starting gRPC Server on :%d", configs.Env.RpcPort)

	app.Serve(l)
}
