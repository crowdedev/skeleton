package dics

import (
	events "github.com/crowdeco/skeleton/events"
	"github.com/sarulabs/dingo/v4"
)

var Dispatcher = []dingo.Def{
	{
		Name: "core:event:dispatcher",
		Build: func() (*events.Dispatcher, error) {
			return events.NewDispatcher([]events.Listener{
				// @see skeleton-todo
			}), nil
		},
	},
}
