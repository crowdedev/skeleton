package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
)

type Database struct {
}

func (d *Database) Run(servers []configs.Server) {
	log.Printf("Starting DB Auto Migration")

	for _, server := range servers {
		server.RegisterAutoMigrate()
	}
}

func (d *Database) IsBackground() bool {
	return true
}
