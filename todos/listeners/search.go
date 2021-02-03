package listeners

import (
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
)

type todoSearch struct {
}

func NewTodoSearch() events.Listener {
	return &todoSearch{}
}

func (s *todoSearch) Listen() string {
	return handlers.PAGINATION_EVENT
}

func (s *todoSearch) Handle(event interface{}) {
	// Example of Listener

	// e := event.(*events.PaginationEvent)
	// query := e.Query()

	// for _, v := range e.Filters() {
	// 	query.Must(elastic.NewTermQuery(v.Field, v.Value))
	// }
}
