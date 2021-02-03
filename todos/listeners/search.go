package listeners

import (
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	adapter "github.com/crowdeco/skeleton/paginations/adapter"
	elastic "github.com/olivere/elastic"
)

type todoSearch struct {
}

func NewTodoSearch() events.Listener {
	return &todoSearch{}
}

func (s *todoSearch) Listen() string {
	return handlers.BEFORE_PAGINATION_EVENT
}

func (s *todoSearch) Handle(event interface{}) {
	e := event.(*adapter.PaginationEvent)
	query := e.Query()

	for _, v := range e.Filters() {
		query.Must(elastic.NewTermQuery(v.Field, v.Value))
	}
}
