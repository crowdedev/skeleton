package configs

import (
	"fmt"

	"github.com/crowdeco/skeleton/configs/driver"

	"gorm.io/gorm"
)

var Database *gorm.DB

func loadDatabase() {
	var db *gorm.DB

	switch Env.DbDriver {
	case "mysql":
		db = driver.NewMySQL().Connect(Env.DbHost, Env.DbPort, Env.DbUser, Env.DbPassword, Env.DbName, Env.Debug)
	case "postgresql":
		db = driver.NewPostgreSQL().Connect(Env.DbHost, Env.DbPort, Env.DbUser, Env.DbPassword, Env.DbName, Env.Debug)
	default:
		panic("Driver not defined")
	}

	Database = db

	fmt.Println("Database configured...")
}
