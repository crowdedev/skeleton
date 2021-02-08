package events

import "github.com/crowdeco/skeleton/configs"

type ModelEvent struct {
	Data    interface{}
	Id      string
	Service configs.Service
}
