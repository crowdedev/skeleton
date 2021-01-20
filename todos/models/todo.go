package models

import (
	configs "github.com/crowdeco/skeleton/configs"
)

type Todo struct {
	configs.Base
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

func (Todo) TableName() string {
	return "todos"
}

func (j Todo) SetCreatedBy(user *configs.User) {
}

func (j Todo) SetUpdatedBy(user *configs.User) {
}

func (j Todo) SetDeletedBy(user *configs.User) {
}

func (Todo) IsSoftDelete() bool {
	return false
}
