package listeners

import (
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
)

type TodoSearch struct {
}

func NewTodoSearch() events.Listener {
	return &TodoSearch{}
}

func (s *TodoSearch) Listen() string {
	return handlers.PAGINATION_EVENT
}

func (s *TodoSearch) Handle(event interface{}) {
	// Example of Listener

	// e := event.(*events.PaginationEvent)
	// query := e.Query()

	// for _, v := range e.Filters() {
	// 	query.Must(elastic.NewTermQuery(v.Field, v.Value))
	// }
}
