package listeners

import (
	"context"

	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Delete struct {
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (u *Delete) Handle(event interface{}) {
	e := event.(*events.ModelEvent)

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", e.Id()))
	result, _ := u.Elasticsearch.Search().Index(e.Service()).Query(query).Do(u.Context)
	for _, hit := range result.Hits.Hits {
		u.Elasticsearch.Delete().Index(e.Service()).Id(hit.Id).Do(u.Context)
	}
}

func (u *Delete) Listen() string {
	return handlers.AFTER_DELETE_EVENT
}
