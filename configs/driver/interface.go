package driver

import "gorm.io/gorm"

type Driver interface {
	Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB
}
