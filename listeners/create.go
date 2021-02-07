package listeners

import (
	"context"
	"encoding/json"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Create struct {
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (c *Create) Handle(event interface{}) {
	e := event.(*events.ModelEvent)

	data, _ := json.Marshal(e.Data())
	c.Elasticsearch.Index().Index(e.Service()).BodyJson(string(data)).Do(c.Context)
}

func (u *Create) Listen() string {
	return handlers.AFTER_CREATE_EVENT
}

func (c *Create) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
