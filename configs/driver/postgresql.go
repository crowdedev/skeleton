package driver

import (
	"fmt"
	"log"

	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSql struct {
}

func (d *PostgreSql) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	var db *gorm.DB
	var err error

	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	if debug {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{})
	} else {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	if err != nil {
		log.Printf("Gorm PostgreSQL: %+v \n", err)
		panic(err)
	}

	return db
}
