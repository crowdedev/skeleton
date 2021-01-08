package models

import (
	"fmt"
	"strconv"

	configs "github.com/crowdeco/skeleton/configs"
)

type Todo struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name;type:varchar(255);not null"  json:"name"`
}

func (j Todo) TableName() string {
	return "todos"
}

func (j Todo) Identifier() string {
	return fmt.Sprintf("%d", j.ID)
}

func (j Todo) SetIdentifier(id string) {
	j.ID, _ = strconv.Atoi(id)
}

func (j Todo) SetCreatedBy(user *configs.User) {
}

func (j Todo) SetUpdatedBy(user *configs.User) {
}

func (j Todo) SetDeletedBy(user *configs.User) {
}

func (j Todo) IsSoftDelete() bool {
	return false
}
