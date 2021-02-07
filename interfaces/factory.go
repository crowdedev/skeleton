package interfaces

import (
	"sort"

	configs "github.com/crowdeco/skeleton/configs"
)

type (
	Application struct {
		Applications []configs.Application
	}
)

func (a *Application) Run(servers []configs.Server) {
	sort.Slice(a.Applications, func(i, j int) bool {
		return a.Applications[i].Priority() > a.Applications[j].Priority()
	})

	for _, application := range a.Applications {
		if application.IsBackground() {
			go application.Run(servers)
		} else {
			application.Run(servers)
		}
	}
}
