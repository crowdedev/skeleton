package dics

import (
	events "github.com/crowdeco/skeleton/events"
	"github.com/sarulabs/dingo/v4"
)

var Dispatcher = []dingo.Def{
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
}
