package modules

import (
	events "github.com/crowdeco/skeleton/events"
	todos "github.com/crowdeco/skeleton/todos"
	listeners "github.com/crowdeco/skeleton/todos/listeners"
	"github.com/sarulabs/dingo/v4"
)

var Todo = []dingo.Def{
	{
		Name: "module:todo",
		Build: func(dispatcher *events.Dispatcher) (todos.TodoModule, error) {
			return todos.NewTodoModule(dispatcher), nil
		},
	},
	{
		Name:  "module:todo:listener:search",
		Build: (*listeners.TodoSearch)(nil),
	},
}
