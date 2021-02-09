package events

import (
	services "github.com/crowdeco/skeleton/services"
)

type Model struct {
	Data       interface{}
	Id         string
	Repository *services.Repository
}
