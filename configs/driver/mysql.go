package driver

import (
	"fmt"
	"log"

	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysql struct {
}

func NewMySQL() Driver {
	return &mysql{}
}

func (d *mysql) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	var db *gorm.DB
	var err error

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", user, password, host, port, dbname)
	if debug {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{})
	} else {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	if err != nil {
		log.Printf("Gorm MySQL: %+v \n", err)
		panic(err)
	}

	return db
}
