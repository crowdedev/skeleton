package creates

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

func (c *Elasticsearch) Handle(event interface{}) {
	e := event.(*events.ModelEvent)

	m := e.Data.(configs.Model)
	data, _ := json.Marshal(e.Data)
	c.Elasticsearch.Index().Index(m.TableName()).BodyJson(string(data)).Do(c.Context)
}

func (u *Elasticsearch) Listen() string {
	return handlers.AFTER_CREATE_EVENT
}

func (c *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
