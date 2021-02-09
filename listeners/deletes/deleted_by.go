package deletes

import (
	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
)

type DeletedBy struct {
	Env *configs.Env
}

func (c *DeletedBy) Handle(event interface{}) {
	e := event.(*events.Model)
	data := e.Data.(configs.Model)
	data.SetDeletedBy(c.Env.User)
	e.Repository.OverrideData(data)
}

func (u *DeletedBy) Listen() string {
	return handlers.BEFORE_DELETE_EVENT
}

func (c *DeletedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
