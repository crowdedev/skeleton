package dics

import (
	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/paginations/adapter"
	"github.com/KejawenLab/bima/v3/repositories"
	"github.com/sarulabs/dingo/v4"
)

var Container = []dingo.Def{
	{
		Name:  "bima:repository:gorm",
		Scope: bima.Application,
		Build: (*repositories.GormRepository)(nil),
	},
	{
		Name:  "bima:pagination:adapter:gorm",
		Scope: bima.Application,
		Build: func(
			env *configs.Env,
			dispatcher *events.Dispatcher,
		) (*adapter.GormAdapter, error) {
			return &adapter.GormAdapter{
				Debug:      env.Debug,
				Dispatcher: dispatcher,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name:  "bima:handler",
		Scope: bima.Application,
		Build: func(
			env *configs.Env,
			dispatcher *events.Dispatcher,
			repository repositories.Repository,
			adapter paginations.Adapter,
		) (*handlers.Handler, error) {
			return &handlers.Handler{
				Debug:      env.Debug,
				Dispatcher: dispatcher,
				Repository: repository,
				Adapter:    adapter,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:event:dispatcher"),
			"2": dingo.Service("bima:repository:gorm"),
			"3": dingo.Service("bima:pagination:adapter:gorm"),
		},
	},
}
