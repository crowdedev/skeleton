package dics

import (
	events "github.com/crowdeco/skeleton/events"
	"github.com/crowdeco/skeleton/interfaces"
	"github.com/sarulabs/dingo/v4"
)

var Core = []dingo.Def{
	{
		Name: "core:event:dispatcher",
		Build: func(
			todo events.Listener,
		) (*events.Dispatcher, error) {
			return events.NewDispatcher([]events.Listener{todo}), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("module:todo:listener:search"),
		},
	},
	{
		Name:  "core:interface:database",
		Build: (*interfaces.Database)(nil),
		Params: dingo.Params{
			"Dispatcher": dingo.Service("core:event:dispatcher"),
		},
	},
	{
		Name:  "core:interface:grpc",
		Build: (*interfaces.GRpc)(nil),
		Params: dingo.Params{
			"Dispatcher": dingo.Service("core:event:dispatcher"),
		},
	},
	{
		Name:  "core:interface:queue",
		Build: (*interfaces.Queue)(nil),
		Params: dingo.Params{
			"Dispatcher": dingo.Service("core:event:dispatcher"),
		},
	},
	{
		Name:  "core:interface:rest",
		Build: (*interfaces.Rest)(nil),
	},
}
