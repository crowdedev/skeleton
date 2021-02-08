package creates

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
)

type CreatedBy struct {
	Env *configs.Env
}

func (c *CreatedBy) Handle(event interface{}) {
	e := event.(*events.ModelEvent)
	data := e.Data.(configs.Model)
	data.SetCreatedBy(c.Env.User)
	e.Service.OverrideData(data)
}

func (u *CreatedBy) Listen() string {
	return handlers.BEFORE_CREATE_EVENT
}

func (c *CreatedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
