package events

import (
	services "github.com/crowdeco/skeleton/services"
)

type ModelEvent struct {
	Data    interface{}
	Id      string
	Service *services.Service
}
