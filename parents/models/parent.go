package models

import (
	configs "github.com/crowdeco/skeleton/configs"
)

type Parent struct {
	configs.Base
	Name       string             `gorm:"column:name;type:varchar(255);notNull"`
	Nullable   configs.NullString `gorm:"column:nullable;type:varchar(255)"`
	CreateOnly string             `gorm:"column:name;type:varchar(255);notNull"`
	Children   []Child            ``
}

func (Parent) TableName() string {
	return "parents"
}

func (j Parent) SetCreatedBy(user *configs.User) {
}

func (j Parent) SetUpdatedBy(user *configs.User) {
}

func (j Parent) SetDeletedBy(user *configs.User) {
}

func (Parent) IsSoftDelete() bool {
	return false
}
