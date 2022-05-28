package dics

import (
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/sarulabs/dingo/v4"
)

var Container = []dingo.Def{
	{
		Name:  "bima:handler:handler",
		Build: (*handlers.Handler)(nil),
		Params: dingo.Params{
			"Env":        dingo.Service("bima:config:env"),
			"Context":    dingo.Service("bima:context:background"),
			"Dispatcher": dingo.Service("bima:event:dispatcher"),
			"Repository": dingo.Service("bima:service:repository"),
			"Adapter":    dingo.Service("bima:pagination:adapter:elasticsearch"),
		},
	},
}
