package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	creates "github.com/crowdeco/skeleton/listeners/creates"
	deletes "github.com/crowdeco/skeleton/listeners/deletes"
	updates "github.com/crowdeco/skeleton/listeners/updates"
	"github.com/sarulabs/dingo/v4"
)

var Dispatcher = []dingo.Def{
	{
		Name: "core:event:dispatcher",
		Build: func(
			a configs.Listener,
			b configs.Listener,
			c configs.Listener,
			d configs.Listener,
			e configs.Listener,
			f configs.Listener,
		) (*events.Dispatcher, error) {
			dispatcher := events.Dispatcher{
				Events: make(map[string][]configs.Listener),
			}

			dispatcher.Register([]configs.Listener{a, b, c, d, e, f})

			return &dispatcher, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:listener:create:elasticsearch"),
			"1": dingo.Service("core:listener:update:elasticsearch"),
			"2": dingo.Service("core:listener:delete:elasticsearch"),
			"3": dingo.Service("core:listener:create:created_by"),
			"4": dingo.Service("core:listener:update:updated_by"),
			"5": dingo.Service("core:listener:delete:deleted_by"),
		},
	},
	{
		Name:  "core:listener:create:elasticsearch",
		Build: (*creates.Elasticsearch)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:update:elasticsearch",
		Build: (*updates.Elasticsearch)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:delete:elasticsearch",
		Build: (*deletes.Elasticsearch)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:create:created_by",
		Build: (*creates.CreatedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
	{
		Name:  "core:listener:update:updated_by",
		Build: (*updates.UpdatedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
	{
		Name:  "core:listener:delete:deleted_by",
		Build: (*deletes.DeletedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
}
