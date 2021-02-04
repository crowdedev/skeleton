package interfaces

import (
	"log"

	configs "github.com/crowdeco/skeleton/configs"
)

type Database struct {
	Servers []configs.Server
}

func (d *Database) Run() {
	log.Printf("Starting DB Auto Migration")

	for _, server := range d.Servers {
		server.RegisterAutoMigrate()
	}
}
