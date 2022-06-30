package dics

import (
	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/repositories"
	"github.com/sarulabs/dingo/v4"
)

var Container = []dingo.Def{
	{
		Name: "bima:handler",
		Build: func(
			env *configs.Env,
			logger *loggers.Logger,
			dispatcher *events.Dispatcher,
			repository repositories.Repository,
			adapter paginations.Adapter,
		) (*handlers.Handler, error) {
			return &handlers.Handler{
				Debug:      env.Debug,
				Logger:     logger,
				Dispatcher: dispatcher,
				Repository: repository,
				Adapter:    adapter,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:logger"),
			"2": dingo.Service("bima:event:dispatcher"),
			"3": dingo.Service("bima:service:repository:gorm"),
			"4": dingo.Service("bima:pagination:adapter:gorm"),
		},
	},
}
