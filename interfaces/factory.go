package interfaces

import (
	"sort"

	"github.com/crowdeco/skeleton/configs"
)

type Application struct {
	Applications []configs.Application
	Servers      []configs.Server
}

func (a *Application) Run() {
	sort.Slice(a.Applications, func(i, j int) bool {
		return a.Applications[i].Priority() > a.Applications[j].Priority()
	})

	for _, application := range a.Applications {
		if application.IsBackground() {
			go application.Run(a.Servers)
		} else {
			application.Run(a.Servers)
		}
	}
}
