package updates

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
)

type UpdatedBy struct {
	Env *configs.Env
}

func (c *UpdatedBy) Handle(event interface{}) {
	e := event.(*events.ModelEvent)
	data := e.Data().(configs.Model)
	data.SetUpdatedBy(c.Env.User)
}

func (u *UpdatedBy) Listen() string {
	return handlers.BEFORE_DELETE_EVENT
}

func (c *UpdatedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}