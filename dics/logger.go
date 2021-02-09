package dics

import (
	configs "github.com/crowdeco/skeleton/configs"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

var Logger = []dingo.Def{
	{
		Name: "core:logger:extension",
		Build: func(
			mongodb logrus.Hook,
		) (*configs.LoggerExtension, error) {
			return &configs.LoggerExtension{
				Extensions: []logrus.Hook{
					mongodb,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:logger:extension:mongodb"),
		},
	},
}
