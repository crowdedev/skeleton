package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
)

type Elasticsearch struct {
}

func (e *Elasticsearch) Run(servers []configs.Server) {
	log.Printf("Repopulating Data")

	for _, server := range servers {
		server.RepopulateData()
	}
}

func (e *Elasticsearch) IsBackground() bool {
	return true
}
