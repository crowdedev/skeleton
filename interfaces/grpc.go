package interfaces

import (
	"fmt"
	"log"
	"net"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	grpc "google.golang.org/grpc"
)

type GRpc struct {
	Env        *configs.Env
	GRpc       *grpc.Server
	Dispatcher *events.Dispatcher
	Servers    []configs.Server
}

func (g *GRpc) Register(servers []configs.Server) {
	g.Servers = servers
}

func (g *GRpc) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Env.RpcPort))
	if err != nil {
		log.Fatalf("Port %d is not available. %v", g.Env.RpcPort, err)
	}

	for _, server := range g.Servers {
		server.RegisterGRpc(g.GRpc)
	}

	log.Printf("Starting gRPC Server on :%d", g.Env.RpcPort)

	g.GRpc.Serve(l)
}
