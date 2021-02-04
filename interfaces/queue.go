package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
)

type Queue struct {
	Servers []configs.Server
}

func (q *Queue) Run() {
	log.Printf("Starting Queue Consumer")

	for _, server := range q.Servers {
		go server.RegisterQueueConsumer()
	}
}
