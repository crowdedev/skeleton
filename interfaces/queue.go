package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
)

type Queue struct {
}

func (q *Queue) Run(servers []configs.Server) {
	log.Printf("Starting Queue Consumer")

	for _, server := range servers {
		go server.RegisterQueueConsumer()
	}
}

func (q *Queue) IsBackground() bool {
	return true
}

func (q *Queue) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
