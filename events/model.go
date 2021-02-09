package events

import (
	services "github.com/crowdeco/skeleton/services"
)

type ModelEvent struct {
	Data       interface{}
	Id         string
	Repository *services.Repository
}
