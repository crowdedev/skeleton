package interfaces

import (
	"github.com/crowdeco/skeleton/configs"
)

type Application struct {
	Applications []configs.Application
	Servers      []configs.Server
}

func (a *Application) Run() {
	for _, application := range a.Applications {
		if application.IsBackground() {
			go application.Run(a.Servers)
		} else {
			application.Run(a.Servers)
		}
	}
}
