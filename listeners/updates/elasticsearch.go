package updates

import (
	"context"
	"encoding/json"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (u *Elasticsearch) Handle(event interface{}) {
	e := event.(*events.Model)

	m := e.Data.(configs.Model)
	query := elastic.NewMatchQuery("ID", e.Id)
	result, _ := u.Elasticsearch.Search().Index(m.TableName()).Query(query).Do(u.Context)
	for _, hit := range result.Hits.Hits {
		u.Elasticsearch.Delete().Index(m.TableName()).Id(hit.Id).Do(u.Context)
	}

	data, _ := json.Marshal(e.Data)
	u.Elasticsearch.Index().Index(m.TableName()).BodyJson(string(data)).Do(u.Context)
}

func (u *Elasticsearch) Listen() string {
	return handlers.AFTER_UPDATE_EVENT
}

func (u *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
