package todos

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/repositories"
	grpcs "github.com/KejawenLab/skeleton/v3/protos/builds"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	*bima.Server
	Module *Module
}

func (s *Server) Register(gs *grpc.Server) {
	grpcs.RegisterTodosServer(gs, s.Module)
}

func (s *Server) Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	return grpcs.RegisterTodosHandler(context, server, client)
}

func (s *Server) Migrate(db *gorm.DB) {
	if s.Debug {
		db.AutoMigrate(&Todo{})
	}
}

func (s *Server) Sync(client *elastic.Client) {
	ctx := context.WithValue(context.Background(), "scope", "sync")
	var records []Todo
	err := s.Module.Handler.Repository.FindBy(&records, repositories.Filter{
		Field:    "synced_at",
		Operator: "<=",
		Value:    time.Now().Add(-5 * time.Minute), //Last sync 5 minutes ago
	})
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())
	}

	index := fmt.Sprintf("%s_%s", s.Module.Model.Env.Service.ConnonicalName, s.Module.Model.TableName())
	for _, d := range records {
		data, _ := json.Marshal(d)
		if d.SyncedAt.Valid {
			query := elastic.NewMatchQuery("Id", d.Id)
			result, _ := client.Search().Index(index).Query(query).Do(ctx)
			if result != nil && result.Hits != nil {
				for _, hit := range result.Hits.Hits {
					client.Delete().Index(index).Id(hit.Id).Do(ctx)
				}
			}

			data, _ := json.Marshal(d)
			client.Index().Index(index).BodyJson(string(data)).Do(ctx)
		} else {
			client.Index().Index(index).BodyJson(string(data)).Do(ctx)
		}

		d.SetSyncedAt(time.Now())
		d.Env = s.Module.Model.Env
		s.Module.Handler.Repository.Update(&d)
	}
}
