package dics

import (
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/sarulabs/dingo/v4"
)

var Container = []dingo.Def{
	{
		Name:  "bima:handler:handler",
		Build: (*handlers.Handler)(nil),
		Params: dingo.Params{
			"Dispatcher": dingo.Service("bima:event:dispatcher"),
			"Repository": dingo.Service("bima:service:repository:gorm"),
			"Adapter":    dingo.Service("bima:pagination:adapter:gorm"),
			"Logger":     dingo.Service("bima:handler:logger"),
		},
	},
}
