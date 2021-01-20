package configs

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func loadDatabase() {
	host := Env.DbHost
	port := Env.DbPort
	user := Env.DbUser
	password := Env.DbPassword
	dbname := Env.DbName

	var db *gorm.DB
	var err error

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", user, password, host, port, dbname)

	if Env.Debug {
		db, err = gorm.Open(mysql.Open(conn), &gorm.Config{})
	} else {
		db, err = gorm.Open(mysql.Open(conn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	if err != nil {
		log.Printf("Gorm: %+v \n", err)
		panic(err)
	}

	Database = db

	fmt.Println("Database configured...")
}
