package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	listeners "github.com/crowdeco/skeleton/listeners"
	"github.com/sarulabs/dingo/v4"
)

var Dispatcher = []dingo.Def{
	{
		Name: "core:event:dispatcher",
		Build: func(
			create configs.Listener,
			update configs.Listener,
			delete configs.Listener,
		) (*events.Dispatcher, error) {
			dispatcher := events.Dispatcher{
				Events: make(map[string]configs.Listener),
			}

			dispatcher.Register([]configs.Listener{
				create,
				update,
				delete,
			})

			return &dispatcher, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:listener:create"),
			"1": dingo.Service("core:listener:update"),
			"2": dingo.Service("core:listener:delete"),
		},
	},
	{
		Name:  "core:listener:create",
		Build: (*listeners.Create)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:update",
		Build: (*listeners.Update)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:delete",
		Build: (*listeners.Delete)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
}
