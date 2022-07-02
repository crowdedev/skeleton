package dics

import (
	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/paginations/adapter"
	"github.com/KejawenLab/bima/v3/repositories"
	"github.com/olivere/elastic/v7"
	"github.com/sarulabs/dingo/v4"
	"gorm.io/gorm"
)

var Container = []dingo.Def{
	{
		Name: "bima:pagination:adapter:gorm",
		Build: func(
			env *configs.Env,
			db *gorm.DB,
			dispatcher *events.Dispatcher,
		) (*adapter.GormAdapter, error) {
			return &adapter.GormAdapter{
				Debug:      env.Debug,
				Database:   db,
				Dispatcher: dispatcher,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:database"),
			"2": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name: "bima:pagination:adapter:elasticsearch",
		Build: func(env *configs.Env, client *elastic.Client, dispatcher *events.Dispatcher) (*adapter.ElasticsearchAdapter, error) {
			return &adapter.ElasticsearchAdapter{
				Debug:      env.Debug,
				Service:    env.Service.ConnonicalName,
				Client:     client,
				Dispatcher: dispatcher,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:elasticsearch:client"),
			"2": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name: "bima:handler",
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
			"3": dingo.Service("bima:pagination:adapter:elasticsearch"),
		},
	},
}
