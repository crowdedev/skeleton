package configs

import "fmt"

func _RegisterAutoMigration() {
	if Env.DbAutoMigrate {
		fmt.Println("Running Auto Migration...")
		Database.AutoMigrate()
	}
}
