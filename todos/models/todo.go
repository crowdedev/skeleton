package models

import (
	"fmt"
	"strconv"
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
