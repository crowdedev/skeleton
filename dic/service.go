package dic

import (
	events "github.com/crowdeco/skeleton/events"
	todos "github.com/crowdeco/skeleton/todos"
	"github.com/sarulabs/dingo/v4"
)

var services = []dingo.Def{
	{
		Name: "core:dispatcher",
		Build: func() (*events.Dispatcher, error) {
			return events.NewDispatcher(), nil
		},
	},
	{
		Name: "module:todo",
		Build: func(dispatcher *events.Dispatcher) (todos.TodoModule, error) {
			return todos.NewTodoModule(dispatcher), nil
		},
	},
}
