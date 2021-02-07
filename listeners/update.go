package listeners

import (
	"context"
	"encoding/json"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Update struct {
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (u *Update) Handle(event interface{}) {
	e := event.(*events.ModelEvent)

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", e.Id()))
	result, _ := u.Elasticsearch.Search().Index(e.Service()).Query(query).Do(u.Context)
	for _, hit := range result.Hits.Hits {
		u.Elasticsearch.Delete().Index(e.Service()).Id(hit.Id).Do(u.Context)
	}

	data, _ := json.Marshal(e.Data())
	u.Elasticsearch.Index().Index(e.Service()).BodyJson(string(data)).Do(u.Context)
}

func (u *Update) Listen() string {
	return handlers.AFTER_UPDATE_EVENT
}

func (u *Update) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
