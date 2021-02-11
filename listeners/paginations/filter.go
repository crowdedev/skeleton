package paginations

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	"github.com/olivere/elastic/v7"
)

type Filter struct {
}

func (u *Filter) Handle(event interface{}) {
	e := event.(*events.Pagination)
	query := e.Query
	filters := e.Filters

	for _, v := range filters {
		query.Must(elastic.NewMatchQuery(v.Field, v.Value))
	}
}

func (u *Filter) Listen() string {
	return handlers.PAGINATION_EVENT
}

func (u *Filter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
