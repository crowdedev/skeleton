package parents

import (
	configs "github.com/crowdeco/skeleton/configs"
	"github.com/crowdeco/skeleton/events"
	models "github.com/crowdeco/skeleton/parents/models"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	"google.golang.org/grpc"
)

type server struct {
	parent ParentModule
	child  ChildModule
}

func NewServer(dispatcher *events.Dispatcher) configs.Server {
	return &server{
		parent: NewParentModule(dispatcher),
		child:  NewChildModule(dispatcher),
	}
}

func (s *server) RegisterGRpc(gs *grpc.Server) {
	grpcs.RegisterParentsServer(gs, s.parent)
	grpcs.RegisterChildrenServer(gs, s.child)
}

func (s *server) RegisterAutoMigrate() {
	if configs.Env.DbAutoMigrate {
		configs.Database.AutoMigrate(&models.Parent{})
		configs.Database.AutoMigrate(&models.Child{})
	}
}

func (s *server) RegisterQueueConsumer() {
	s.parent.Consume()
	s.child.Consume()
}
