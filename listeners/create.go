package listeners

import (
	"context"
	"encoding/json"

	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Create struct {
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (u *Create) Handle(event interface{}) {
	e := event.(*events.ModelEvent)

	data, _ := json.Marshal(e.Data())
	u.Elasticsearch.Index().Index(e.Service()).BodyJson(string(data)).Do(u.Context)
}

func (u *Create) Listen() string {
	return handlers.AFTER_CREATE_EVENT
}
