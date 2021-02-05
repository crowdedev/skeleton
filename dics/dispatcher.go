package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	"github.com/sarulabs/dingo/v4"
)

var Dispatcher = []dingo.Def{
	{
		Name: "core:event:dispatcher",
		Build: func() (*events.Dispatcher, error) {
			dispatcher := events.Dispatcher{
				Events: make(map[string]configs.Listener),
			}

			dispatcher.Register([]configs.Listener{
				// @see skeleton-todo
			})

			return &dispatcher, nil
		},
	},
}
